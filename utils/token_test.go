package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTokenPayload(t *testing.T) {
	username := RandomUsername()

	payload, err := NewTokenPayload(username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiresAt, time.Now().Add(time.Minute), time.Second)

	err = payload.Valid()
	require.NoError(t, err)
}

func TestInvalidTokenPayload(t *testing.T) {
	username := RandomUsername()

	payload, err := NewTokenPayload(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiresAt, time.Now().Add(-time.Minute), time.Second)

	err = payload.Valid()
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestTokenMaker(t *testing.T) {
	secret := "12345678901234567890123456789012"
	maker, err := NewTokenMaker(secret)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := RandomUsername()
	token, err := maker.Create(username, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.Validate(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Username, username)

	err = payload.Valid()
	require.NoError(t, err)
}

func TestInvalidTokenSecret(t *testing.T) {
	secret := "secret"
	maker, err := NewTokenMaker(secret)
	require.EqualError(t, err, ErrInvalidTokenSecret.Error())
	require.Empty(t, maker)
}

func TestExpiredToken(t *testing.T) {
	secret := "12345678901234567890123456789012"
	maker, err := NewTokenMaker(secret)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := RandomUsername()
	token, err := maker.Create(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	_, err = maker.Validate(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
}

func TestInvalidToken(t *testing.T) {
	secret := "12345678901234567890123456789012"
	maker, err := NewTokenMaker(secret)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	_, err = maker.Validate("token")
	require.EqualError(t, err, ErrInvalidToken.Error())
}
