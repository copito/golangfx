package auth

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type BasicAuthProvider struct {
	logger       *slog.Logger
	jwtValidator *JWTDecoder
	decodedToken *jwt.MapClaims
}

func NewBasicAuthProvider(logger *slog.Logger, jwtValidator *JWTDecoder) *BasicAuthProvider {
	if jwtValidator == nil {
		jwtValidator = NewJWTDecoder("https://auth.example.com", "token/publickey", "")
	}

	return &BasicAuthProvider{
		logger:       logger,
		jwtValidator: jwtValidator,
	}
}

func (v *BasicAuthProvider) ParseToken(token string) error {
	if token == "" {
		return errors.New("invalid Bearer token provided")
	}

	decoded, err := v.jwtValidator.DecodeToken(token)
	if err != nil {
		return err
	}

	v.decodedToken = &decoded
	return nil
}

func (v *BasicAuthProvider) IsValid(token string) bool {
	if v.decodedToken != nil {
		return true
	}

	if token == "" {
		return false
	}

	err := v.ParseToken(token)
	return err == nil
}

func (v *BasicAuthProvider) GetUser(token string) (string, error) {
	if v.decodedToken != nil {
		user, err := v.decodedToken.GetSubject()
		if err != nil {
			return "", err
		}

		// Replace the prefix "KEYCLOCK\\"
		username := strings.Replace(user, "KEYCLOCK\\", "", 1)
		return username, nil
	}

	if token == "" {
		return "", errors.New("invalid Bearer token provided")
	}

	err := v.ParseToken(token)
	if err != nil {
		return "", err
	}

	user, err := v.decodedToken.GetSubject()
	if err != nil {
		return "", err
	}

	// Replace the prefix "KEYCLOCK\\"
	username := strings.Replace(user, "KEYCLOCK\\", "", 1)

	// sub = decoded["sub"].replace("KEYCLOCK\\", "")
	// return {"username": sub}
	return username, nil
}
