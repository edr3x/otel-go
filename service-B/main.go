package main

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/edr3x/otel-go/pkg/entities"
	"github.com/edr3x/otel-go/pkg/entities/responders"
	"github.com/edr3x/otel-go/pkg/otelx"
	"github.com/edr3x/otel-go/service-b/internal/handlers"
)

func init() {
	zap.ReplaceGlobals(createProductionLogger())
}

func main() {
	// for demo purpose set env
	os.Setenv("ENV", "development")
	os.Setenv("OTEL_ENABLE", "true")
	os.Setenv("SERVICE_NAME", "Second test Service")
	os.Setenv("OTEL_COLLECTOR_ENDPOINT", "localhost:4317")

	serviceName := os.Getenv("SERVICE_NAME")

	ctx := context.Background()

	tp, shutdown := otelx.NewTraceProvider(ctx, serviceName)
	defer shutdown()

	shutdownMeterProvider := otelx.NewMeterProvider(ctx, serviceName)
	defer shutdownMeterProvider()

	r := echo.New()
	r.HTTPErrorHandler = entities.CentralEchoErrorHandler

	r.Use(middleware.RequestID())
	r.Use(otelecho.Middleware(serviceName, otelecho.WithTracerProvider(tp))) // otelecho middleware
	r.Use(echo.WrapMiddleware(otelx.MetricsMiddleware))                      // otel metrics
	r.Use(entities.CustomRequestLoggerConfig())                              // must be after otel middleware to extract traceID

	r.GET("/", func(c echo.Context) error {
		return c.String(200, "hello there")
	})

	responder := responders.NewResponder()
	h := handlers.NewHandler(responder)

	r.GET("/posts/:id", h.Foo)

	r.RouteNotFound("/*", func(c echo.Context) error {
		return entities.ErrorNotFound("Route Not Found")
	})

	http.ListenAndServe("0.0.0.0:8081", r)
}

func createProductionLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	return zap.Must(config.Build())
}
