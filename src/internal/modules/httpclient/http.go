package httpclient

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/modules/config"
)

type Params struct {
	fx.In
	Logger         *slog.Logger
	ConfigProvider config.ConfigProvider
}

type Result struct {
	fx.Out

	Client *http.Client
}

func sanitizeURL(url string) string {
	// Simple example: replace digits with 'id'
	// You can enhance this function based on your specific requirements.

	sanitized := url

	// Replace any numbers using regex \d+ with <sanitized>
	// For simplicity, let's just replace any sequence of digits with "<sanitized>"
	regexDigits := regexp.MustCompile(`\d+`)
	sanitized = regexDigits.ReplaceAllString(sanitized, "<sanitized>")

	// Replace UUIDs using regex with <sanitized>
	regexUUID := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	sanitized = regexUUID.ReplaceAllString(sanitized, "<sanitized>")

	// Replace email addresses using regex with <sanitized>
	regexEmail := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	sanitized = regexEmail.ReplaceAllString(sanitized, "<sanitized>")

	return sanitized
}

func NewHTTPClient(params Params) (Result, error) {
	params.Logger.Info("setting up HTTP Client module...")

	transport := otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			sanitizedURL := sanitizeURL(r.URL.Path)
			return fmt.Sprintf("%s %s", r.Method, sanitizedURL)
		}),
	)

	client := &http.Client{
		Transport: transport,
	}

	return Result{Client: client}, nil
}
