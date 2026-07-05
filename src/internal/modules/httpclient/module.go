package httpclient

import "go.uber.org/fx"

var Module = fx.Provide(NewHTTPClient)
