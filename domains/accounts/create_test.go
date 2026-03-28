package accounts

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MaraSystems/graybank_api/api"
	mockdb "github.com/MaraSystems/graybank_api/db/mock"
	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/middlewares"
	"github.com/MaraSystems/graybank_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	account := DummyAccount()

	testCases := []struct {
		name          string
		body          gin.H
		authorize     func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker)
		buildStub     func(t *testing.T, store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"currency": "NGN",
			},
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				middlewares.AddDummyAuthorization(t, request, tokenMaker, "Bearer", "user", time.Minute)
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusCreated)
			},
		},
		{
			name: "UnAuthorized",
			body: gin.H{
				"currency": "NGN",
			},
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "CurrencyConflict",
			body: gin.H{
				"currency": "NGN",
			},
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				middlewares.AddDummyAuthorization(t, request, tokenMaker, "Bearer", "user", time.Minute)
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, &pgconn.PgError{ConstraintName: "owner_currency_key"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusForbidden)
			},
		},
		{
			name: "CurrencyInvalid",
			body: gin.H{
				"currency": "NG",
			},
			authorize: func(t *testing.T, request *http.Request, tokenMaker *utils.TokenMaker) {
				middlewares.AddDummyAuthorization(t, request, tokenMaker, "Bearer", "user", time.Minute)
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStub(t, store)

			server, err := api.TestServer(t, store)
			require.NoError(t, err)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			SetAccountsRoutes(server)
			tc.authorize(t, request, server.TokenMaker)

			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
