package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	ErrExpiredToken       = fmt.Errorf("token has expired")
	ErrInvalidToken       = fmt.Errorf("token is invalid")
	ErrInvalidTokenSecret = fmt.Errorf("invalid secret; must have %d characters", chacha20poly1305.KeySize)
)

type TokenPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewTokenPayload(username string, duration time.Duration) (payload *TokenPayload, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}

	payload = &TokenPayload{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	return
}

func (t *TokenPayload) Valid() error {
	if time.Now().After(t.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}

type TokenMaker struct {
	secretKey []byte
	engine    *paseto.V2
}

func NewTokenMaker(secret string) (maker *TokenMaker, err error) {
	if len(secret) != chacha20poly1305.KeySize {
		err = ErrInvalidTokenSecret
		return
	}

	maker = &TokenMaker{
		secretKey: []byte(secret),
		engine:    paseto.NewV2(),
	}
	return
}

func (t *TokenMaker) Create(username string, duration time.Duration) (token string, err error) {
	payload, err := NewTokenPayload(username, duration)
	if err != nil {
		return
	}

	return t.engine.Encrypt(t.secretKey, payload, nil)
}

func (t *TokenMaker) Validate(token string) (payload *TokenPayload, err error) {
	payload = &TokenPayload{}
	err = t.engine.Decrypt(token, t.secretKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	return
}
