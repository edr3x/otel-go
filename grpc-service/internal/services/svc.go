package services

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/edr3x/otel-go/grpc-service/pkg/pb/proto"
	"github.com/edr3x/otel-go/pkg/otelx"
)

type AssetServer struct {
	proto.AssetServiceServer
}

func (s *AssetServer) GetAssetById(ctx context.Context, req *proto.GetAssetRequest) (*proto.GetAssetResponse, error) {
	dbres, err := getAssetQuery(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	res := &proto.GetAssetResponse{
		Id:      req.Id,
		Key:     dbres.Key,
		Url:     dbres.URL,
		AltText: dbres.AltText,
	}
	return res, nil
}

type assetRes struct {
	Key     string
	URL     string
	AltText string
}

func getAssetQuery(ctx context.Context, uid string) (assetRes, error) {
	_, span := otelx.StartSpan(ctx)
	defer span.End()

	time.Sleep(3 * time.Second)

	res := assetRes{}

	if uid == "9999" {
		return res, status.Error(codes.InvalidArgument, "the provided uuid is invalid")
	}

	res.Key = "orange-cat.jpg"
	res.URL = "https://www.alleycat.org/wp-content/uploads/2019/03/FELV-cat.jpg"
	res.AltText = "Orange Cat "

	return res, nil
}
