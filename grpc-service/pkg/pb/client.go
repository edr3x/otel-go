package pb

import (
	"os"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/edr3x/otel-go/grpc-service/pkg/pb/proto"
)

type AssetClient struct {
	proto.AssetServiceClient
	conn *grpc.ClientConn
}

var client *AssetClient

func NewAssetClient() (proto.AssetServiceClient, error) {
	if client != nil && client.conn != nil {
		return client.AssetServiceClient, nil
	}

	// Lock to prevent race conditions
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	// Double-check to avoid multiple initializations
	if client != nil && client.conn != nil {
		return client.AssetServiceClient, nil
	}
	serviceURL, ok := os.LookupEnv("ASSET_SERVICE_URL")
	if !ok {
		serviceURL = "localhost:50051"
	}
	conn, err := grpc.NewClient(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		return nil, err
	}
	client = &AssetClient{
		AssetServiceClient: proto.NewAssetServiceClient(conn),
		conn:               conn,
	}
	return client.AssetServiceClient, nil
}

func Close() {
	if client != nil && client.conn != nil {
		client.conn.Close()
	}
}
