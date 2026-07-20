package chucknorris

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/modules/config"
	customlogger "github.com/copito/runner/src/pkg/logger"
)

type ChuckNorrisGateway interface {
	GetRandomJoke(ctx context.Context, category string) (*JokeResult, error)
}

type Params struct {
	fx.In

	ConfigProvider config.ConfigProvider
	HTTPClient     *http.Client
}

type Result struct {
	ChuckNorrisGateway ChuckNorrisGateway
}

type chuckNorrisGateway struct {
	httpClient *http.Client
	baseURL    url.URL
}

func NewChuckNorrisGateway(params Params) (Result, error) {
	config := params.ConfigProvider.Get()
	baseURL := config.ChuckNorrisGateway.BaseURL

	httpClient := params.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return Result{
		ChuckNorrisGateway: &chuckNorrisGateway{
			httpClient: httpClient,
			baseURL: url.URL{
				Host: baseURL,
			},
		},
	}, nil
}

func (c *chuckNorrisGateway) GetRandomJoke(ctx context.Context, category string) (*JokeResult, error) {
	logger := customlogger.LoggerFromContext(ctx, slog.Default())

	endpoint, err := url.JoinPath(c.baseURL.Host, "/jokes/random")
	if err != nil {
		return nil, errors.New("unable to build path - internal error")
	}

	// Add 2 second timeout
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*2)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, errors.New("error building request")
	}
	req.Header.Set("Context-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("unexpected status code", slog.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result JokeResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error("failed to decode response", slog.Any("error", err))
		return nil, err
	}

	return &result, nil
}
