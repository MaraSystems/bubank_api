package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaraSystems/bubank_api/api"
	mockdb "github.com/MaraSystems/bubank_api/db/mock"
	db "github.com/MaraSystems/bubank_api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	user, password := dummyUser(t)

	testCases := []struct {
		name          string
		body          *gin.H
		buildStub     func(t *testing.T, store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: &gin.H{
				"username": user.Username,
				"Password": password,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				fmt.Println(recorder)
				require.Equal(t, recorder.Code, http.StatusCreated)
			},
		},
		{
			name: "NoUsernameProvided",
			body: &gin.H{
				"Password": password,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				fmt.Println(recorder)
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "NoPasswordProvided",
			body: &gin.H{
				"username": user.Username,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				fmt.Println(recorder)
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "NotFound",
			body: &gin.H{
				"username": user.Username,
				"Password": password,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				fmt.Println(recorder)
				require.Equal(t, recorder.Code, http.StatusNotFound)
			},
		},
		{
			name: "InternalError",
			body: &gin.H{
				"username": user.Username,
				"Password": password,
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Times(1).
					Return(db.User{}, pgx.ErrTxClosed)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				fmt.Println(recorder)
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
		{
			name: "WrongPassword",
			body: &gin.H{
				"username": user.Username,
				"Password": "password",
			},
			buildStub: func(t *testing.T, store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), user.Username).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				fmt.Println(recorder)
				require.Equal(t, recorder.Code, http.StatusForbidden)
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

			url := "/auth"
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			SetAuthRoutes(server)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, *recorder)
		})
	}
}
