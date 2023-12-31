//go:build e2e

package opentelemetry

import (
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	//tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	//builder := MiddlewareBuilder{
	//	Tracer: tracer,
	//}
	//server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	//
}
