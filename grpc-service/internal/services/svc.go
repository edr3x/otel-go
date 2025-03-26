package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/edr3x/otel-go/grpc-service/pkg/pb/proto"
	"github.com/edr3x/otel-go/pkg/otelx"
)

type AssetServer struct {
	proto.AssetServiceServer
}

func (s *AssetServer) GetAssetById(ctx context.Context, req *proto.GetAssetRequest) (*proto.GetAssetResponse, error) {
	uid, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dbres, err := getAssetQuery(ctx, uid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "asset not found")
	}

	res := &proto.GetAssetResponse{
		Id:      uid.String(),
		Key:     dbres.Key,
		Url:     dbres.Url,
		AltText: dbres.AltText,
	}
	return res, nil
}

type assetRes struct {
	Key     string
	Url     string
	AltText string
}

func getAssetQuery(ctx context.Context, uid uuid.UUID) (assetRes, error) {
	_, span := otelx.StartSpan(ctx)
	defer span.End()

	time.Sleep(3 * time.Second)

	res := assetRes{}

	if uid.String() == "0195d2b8-bbbb-72a5-bcb1-27226d04b0c6" {
		return res, status.Error(codes.AlreadyExists, "the provided uuid is already provided")
	}

	res.Key = "orange-cat.jpg"
	res.Url = "https://www.alleycat.org/wp-content/uploads/2019/03/FELV-cat.jpg"
	res.AltText = "Orange Cat "

	return res, nil
}
