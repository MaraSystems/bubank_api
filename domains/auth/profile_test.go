package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MaraSystems/graybank_api/api"
	mockdb "github.com/MaraSystems/graybank_api/db/mock"
	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/middlewares"
	"github.com/MaraSystems/graybank_api/utils"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestGetProfile(t *testing.T) {
	user, _ := DummyUser(t)

	testCases := []struct {
		name          string
		authorize     func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker)
		buildStub     func(t *testing.T, store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				middlewares.AddDummyAuthorization(t, request, tokenMaker, "Bearer", "user", time.Minute)
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "NotProfile",
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				middlewares.AddDummyAuthorization(t, request, tokenMaker, "Bearer", "user", time.Minute)
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusNotFound)
			},
		},
		{
			name: "ServerError",
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				middlewares.AddDummyAuthorization(t, request, tokenMaker, "Bearer", "user", time.Minute)
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
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

			tc.buildStub(t, store)
			SetAuthRoutes(server)

			url := "/auth"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.authorize(t, request, server.TokenMaker)

			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
