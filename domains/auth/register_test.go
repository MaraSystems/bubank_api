package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaraSystems/graybank_api/api"
	mockdb "github.com/MaraSystems/graybank_api/db/mock"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	user, password := DummyUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStub     func(t *testing.T, store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"email":     user.Email,
				"full_name": user.FullName,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusCreated)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  "#user",
				"password":  password,
				"email":     user.Email,
				"full_name": user.FullName,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"email":     "email",
				"full_name": user.FullName,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "DuplicateEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"email":     user.Email,
				"full_name": user.FullName,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, &pgconn.PgError{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusConflict)
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

			url := "/auth/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			SetAuthRoutes(server)
			recorder := httptest.NewRecorder()
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
