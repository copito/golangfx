package mainhandler

import (
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"

	"github.com/copito/runner/src/internal/handler/common"
	"github.com/copito/runner/src/internal/modules/config"
)

type MainHandler interface {
	common.HttpHandlerInterface
}

// MainHandler is an http.Handler that copies its request body
// back to the response.
type mainHandler struct {
	configProvider config.ConfigProvider
}

// NewMainHandler builds a new MainHandler.
func NewMainHandler(configProvider config.ConfigProvider) *mainHandler {
	return &mainHandler{
		configProvider: configProvider,
	}
}

func (h *mainHandler) Pattern() runtime.Pattern {
	// "/"
	pattern, err := runtime.NewPattern(
		common.ValidPatternVersion,
		[]int{
			int(utilities.OpLitPush), 0, // runtime.OpLitPush → Push the literal "docs" (matches /openapi exactly).
			int(utilities.OpPushM), 1, // runtime.OpPushM → Matches a deep wildcard ({filepath:.+}) capturing everything after /docs/.
		},
		[]string{""},
		"", // no verb (gRPC routing suffix)
	)
	if err != nil {
		panic("error registering pattern for swagger handler")
	}
	return pattern
}

func (h *mainHandler) Method() string {
	return "GET"
}

// ServeHTTP handles an HTTP request to the /docs endpoint.
func (h *mainHandler) ServeHTTP() runtime.HandlerFunc {
	type Data struct {
		ProjectName string `json:"project_name"`
		Version     string `json:"version"`
		Environment string `json:"environment"`
	}

	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		// Add http route to show project
		config := h.configProvider.Get()
		data := Data{
			ProjectName: config.Global.Service,
			Version:     config.Global.Version,
			Environment: string(config.Backend.Environment),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}
