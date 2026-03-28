package middlewares

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/MaraSystems/graybank_api/utils"
	"github.com/stretchr/testify/require"
)

func AddDummyAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker *utils.TokenMaker,
	authType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.Create(username, duration)
	require.NoError(t, err)

	authorization := fmt.Sprintf("%s %s", authType, token)
	request.Header.Add(AuthHeaderKey, authorization)
}
