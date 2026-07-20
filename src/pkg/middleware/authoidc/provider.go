package authoidc

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/copito/runner/src/internal/modules/config"
)

type OIDCProvider struct {
	logger   *slog.Logger
	provider *oidc.Provider
	config   oauth2.Config
	verifier *oidc.IDTokenVerifier

	token *oidc.IDToken
}

func NewOIDCProvider(logger *slog.Logger, configProvider config.ConfigProvider) (*OIDCProvider, error) {
	logger.Info("Initializing OIDC Provider")

	config := configProvider.Get()

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, config.Auth.IssuerURL)
	if err != nil {
		logger.Error("unable to connect to OIDC Provider", slog.String("provider", "keycloak"))
		return nil, err
	}

	logger.Info("Managed to connect to OIDC Provider", slog.String("provider", "keycloak"), slog.Any("provider", provider))
	scopes := []string{oidc.ScopeOpenID}
	scopes = append(scopes, config.Auth.Scopes...)

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     config.Auth.ClientID,
		ClientSecret: config.Auth.ClientSecret,
		RedirectURL:  config.Auth.RedirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: scopes,
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.Auth.ClientID})

	return &OIDCProvider{
		logger:   logger,
		provider: provider,
		config:   oauth2Config,
		verifier: verifier,
	}, nil
}

func (v *OIDCProvider) ParseToken(token string) error {
	if token == "" {
		return errors.New("invalid Bearer token provided")
	}

	decoded, err := v.verifier.Verify(context.Background(), token)
	if err != nil {
		return err
	}

	v.token = decoded
	return nil
}

func (v *OIDCProvider) IsValid(token string) bool {
	if v.token != nil {
		return true
	}

	if token == "" {
		return false
	}

	err := v.ParseToken(token)
	return err == nil
}

func (v *OIDCProvider) GetUser(token string) (string, error) {
	if v.token != nil {
		user := v.token.Subject

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

	user := v.token.Subject

	// Replace the prefix "KEYCLOCK\\"
	username := strings.Replace(user, "KEYCLOCK\\", "", 1)

	// sub = decoded["sub"].replace("KEYCLOCK\\", "")
	// return {"username": sub}
	return username, nil
}
