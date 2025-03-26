package main

import (
	"cmp"
	"context"
	"log"
	"net"
	"os"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/attribute"
	otelcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/edr3x/otel-go/grpc-service/internal/services"
	"github.com/edr3x/otel-go/grpc-service/pkg/pb/proto"

	"github.com/edr3x/otel-go/pkg/otelx"
)

func init() {
	zap.ReplaceGlobals(createProductionLogger())
}

func main() {
	os.Setenv("SERVICE_NAME", "Grpc Test Service")
	os.Setenv("OTLP_ENDPOINT", "localhost:4318")

	serviceName := os.Getenv("SERVICE_NAME")

	tp := otelx.NewTraceProvider(serviceName)
	defer func() { _ = tp.Shutdown(context.Background()) }()

	port := cmp.Or(os.Getenv("PORT"), ":50051")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}

	grpcPanicRecoveryHandler := func(p any) (err error) {
		zap.L().Error("recovered from panic", zap.Any("panic", p), zap.ByteString("stac", debug.Stack()))
		return status.Errorf(codes.Internal, "%s", p)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestLoggerInterceptor,
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(grpcPanicRecoveryHandler)),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler(otelgrpc.WithTracerProvider(tp))),
	)

	// healthcheck
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", hpb.HealthCheckResponse_SERVING)
	hpb.RegisterHealthServer(s, healthServer)

	proto.RegisterAssetServiceServer(s, &services.AssetServer{})

	// enable server reflection
	reflection.Register(s)

	zap.L().Info("server started at", zap.Any("address", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %s: %v", port, err)
	}
}

func requestLoggerInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	span := trace.SpanFromContext(ctx)
	resp, err := handler(ctx, req)
	if err != nil {
		span.SetStatus(otelcodes.Error, err.Error())

		st, ok := status.FromError(err)
		if !ok {
			st = status.New(codes.Internal, err.Error())
		}

		span.SetAttributes(
			attribute.String("rpc.error", st.Err().Error()),
		)

		zap.L().Error(
			"Error handling request",
			zap.Error(err),
			zap.String("Status", st.Code().String()),
			zap.String("Method", info.FullMethod),
			zap.String("Trace-id", span.SpanContext().TraceID().String()),
		)

		return nil, st.Err()
	}
	zap.L().Info("Successfully handled request", zap.String("method", info.FullMethod), zap.String("Trace-id", span.SpanContext().TraceID().String()))
	return resp, nil
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
