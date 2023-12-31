package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	web "webdemo/code"
)

const instrumentationName = "webdemo/code/web"

type MiddlewareBuiler struct {
	Tracer trace.Tracer
}

func (m MiddlewareBuiler) Build() web.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			//尝试和客户端的trace结合
			otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			_, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()
			//defer func() {
			//	span.SetName(ctx.MatchedRoute)
			//	span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
			//	span.End()
			//}()
			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))

			//继续加

			//直接调用下一步
			next(ctx)
			span.SetName(ctx.MatchedRoute)
			//把响应码加上去
			span.SetAttributes(attribute.Int("http:status", ctx.RespStatusCode))
		}
	}
}
