package entities

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type failedResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

func CentralEchoErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return // Prevent duplicate JSON responses
	}

	var message string
	var errPayload any
	runtimeEnv := os.Getenv("RUNTIME_ENV")

	code := http.StatusInternalServerError
	message = http.StatusText(code)
	requestID := c.Response().Header().Get("X-Request-Id")

	// Extract span from request context
	span := trace.SpanFromContext(c.Request().Context())

	// Add tracing attributes
	span.SetAttributes(
		attribute.String("request.id", requestID),
		attribute.String("runtime.env", runtimeEnv),
	)
	span.SetStatus(codes.Error, err.Error())

	// Handle specific errors
	switch e := err.(type) {
	case HttpError:
		zap.L().Error(
			err.Error(),
			zap.Int("Status", e.Code),
			zap.String("Caller", e.Caller),
			zap.String("Request-id", requestID),
			zap.String("Method", c.Request().Method),
			zap.String("URI", c.Request().URL.RequestURI()),
			zap.String("Trace-id", span.SpanContext().TraceID().String()),
		)

		span.SetAttributes(
			attribute.String("error.caller", e.Caller),
		)

		if e.Code < 500 {
			code = e.Code
			message = e.Error()
			if _, ok := e.Message.(string); !ok {
				if _, ok2 := e.Message.(error); !ok2 {
					errPayload = e.Message
				}
			}
		}
	default:
		zap.L().Error(
			err.Error(),
			zap.String("Method", c.Request().Method),
			zap.String("URI", c.Request().URL.RequestURI()),
			zap.String("request-id", requestID),
		)
	}

	// https://github.com/labstack/echo/issues/608
	if c.Request().Method == http.MethodHead {
		err = c.NoContent(code)
	} else {
		msg := failedResponse{Success: false, Message: message, Error: errPayload}
		err = c.JSON(code, msg)
	}
	if err != nil {
		zap.L().Error("Custom Error Handler:", zap.Error(err))
	}
}

// Use this instead of request logger
func CustomRequestLoggerConfig() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: false,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				zap.L().Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
					zap.String("Trace-id", trace.SpanFromContext(c.Request().Context()).SpanContext().TraceID().String()),
				)
			}
			return nil
		},
	})
}
