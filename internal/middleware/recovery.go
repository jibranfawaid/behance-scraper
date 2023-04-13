package middleware

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"optimus/internal/utilities"
)

func Recovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, sp := otel.Tracer("middleware").Start(c.Request().Context(), "middleware-process")
			defer sp.End()

			traceId := utilities.Ternary(sp.IsRecording(), sp.SpanContext().TraceID().String(), "(N/A)").(string)
			ctx = context.WithValue(ctx, "traceid", traceId)

			*c.Request() = *(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
