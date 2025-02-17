package handler

import (
	"net/http"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
)

// SwaggerFileHandler is an http.Handler that copies its request body
// back to the response.
type SwaggerFileHandler struct{}

// NewSwaggerFileHandler builds a new SwaggerFileHandler.
func NewSwaggerFileHandler() *SwaggerFileHandler {
	return &SwaggerFileHandler{}
}

func (h *SwaggerFileHandler) Pattern() runtime.Pattern {
	// "/openapi/{filepath:.+}"
	pattern, err := runtime.NewPattern(
		validPatternVersion,
		[]int{
			int(utilities.OpLitPush), 0, // runtime.OpLitPush → Push the literal "openapi" (matches /openapi exactly).
			int(utilities.OpPushM), 1, // runtime.OpPushM → Matches a deep wildcard ({filepath:.+}) capturing everything after /openapi/.
		},
		[]string{"openapi"},
		"", // no verb (gRPC routing suffix)
	)
	if err != nil {
		panic("error registering pattern for swagger file handler")
	}
	return pattern
}

func (h *SwaggerFileHandler) Method() string {
	return "GET"
}

// ServeHTTP handles an HTTP request to the /openapi/* endpoint.
func (h *SwaggerFileHandler) ServeHTTP() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Resolve the correct absolute path to the openapi directory
		openapiDir, err := filepath.Abs("../../../openapi")
		if err != nil {
			http.Error(w, "cannot find openapi path", 400)
		}

		// Serve OpenAPI JSON file properly
		dir := http.Dir(openapiDir)
		fileServer := http.StripPrefix("/openapi/", http.FileServer(dir))
		fileServer.ServeHTTP(w, r)
	}
}
