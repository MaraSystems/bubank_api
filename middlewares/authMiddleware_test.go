package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MaraSystems/bubank_api/api"
	mockdb "github.com/MaraSystems/bubank_api/db/mock"
	"github.com/MaraSystems/bubank_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		authorize     func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				AddDummyAuthorization(t, request, tokenMaker, "Bearer", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "NoAuthorization",
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "InvalidAuthorization",
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				AddDummyAuthorization(t, request, tokenMaker, "", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "InvalidAuthorizationType",
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				AddDummyAuthorization(t, request, tokenMaker, "invalid", "user", time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "ExpiredAuthorizationToken",
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				AddDummyAuthorization(t, request, tokenMaker, "Bearer", "user", -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server, err := api.TestServer(t, store)
			require.NoError(t, err)

			url := "/auth"
			server.Router.GET(
				url,
				AuthMiddleWare(server.TokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.authorize(t, request, server.TokenMaker)
			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
