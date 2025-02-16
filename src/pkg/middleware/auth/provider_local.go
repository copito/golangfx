package auth

import (
	"log/slog"
)

type LocalAuthProvider struct {
	logger *slog.Logger
}

func NewLocalAuthProvider(logger *slog.Logger) *LocalAuthProvider {
	return &LocalAuthProvider{logger: logger}
}

func (v *LocalAuthProvider) IsValid(token string) bool {
	return true
}

func (v *LocalAuthProvider) GetUser(token string) (string, error) {
	return "copito", nil
}
