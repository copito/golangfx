package handler

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	swagger "github.com/swaggest/swgui/v5emb"
)

// SwaggerHandler is an http.Handler that copies its request body
// back to the response.
type SwaggerHandler struct{}

// NewSwaggerHandler builds a new SwaggerHandler.
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

func (h *SwaggerHandler) Pattern() runtime.Pattern {
	// "/docs"
	pattern, err := runtime.NewPattern(
		validPatternVersion,
		[]int{
			int(utilities.OpLitPush), 0, // runtime.OpLitPush → Push the literal "docs" (matches /openapi exactly).
			int(utilities.OpPushM), 1, // runtime.OpPushM → Matches a deep wildcard ({filepath:.+}) capturing everything after /docs/.
		},
		[]string{"docs"},
		"", // no verb (gRPC routing suffix)
	)
	if err != nil {
		panic("error registering pattern for swagger handler")
	}
	return pattern
}

func (h *SwaggerHandler) Method() string {
	return "GET"
}

// ServeHTTP handles an HTTP request to the /docs endpoint.
func (h *SwaggerHandler) ServeHTTP() runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Add swagger as route for http
		// All openapi's are combined into one called services.swagger.json
		// sw := swagger.New("gRPC Gateway API", "/openapi/runner/v1/query.swagger.json", "/docs")
		sw := swagger.New("gRPC Gateway API", "/openapi/services.swagger.json", "/docs")
		sw.ServeHTTP(w, r)
	}
}
