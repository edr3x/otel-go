package handlers

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"

	"github.com/edr3x/otel-go/pkg/entities"
	"github.com/edr3x/otel-go/pkg/entities/responders"
	"github.com/edr3x/otel-go/pkg/otelx"

	assetSvc "github.com/edr3x/otel-go/grpc-service/pkg/pb"
	"github.com/edr3x/otel-go/grpc-service/pkg/pb/proto"
)

type Handler struct {
	res responders.Responders
}

func NewHandler(res responders.Responders) *Handler {
	return &Handler{res: res}
}

func (h *Handler) Foo(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	name, err := getPost(ctx, id)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	return h.res.JSON(c, name)
}

func getPost(ctx context.Context, id string) (*Post, error) {
	ctx, span := otelx.StartSpan(ctx)
	defer span.End()

	time.Sleep(2 * time.Second)

	return getPostByIdQuery(ctx, id)
}

type Post struct {
	Id    string                  `json:"id"`
	Title string                  `json:"title"`
	Asset *proto.GetAssetResponse `json:"asset"`
}

func getPostByIdQuery(ctx context.Context, id string) (*Post, error) {
	_, span := otelx.StartSpan(ctx)
	defer span.End()

	client, err := assetSvc.NewAssetClient()
	if err != nil {
		return nil, entities.ErrorInternal(err)
	}

	if id == "666" {
		span.SetAttributes(attribute.String("id", "666"), attribute.String("extrainfo", "cannot find book of the beast"))
		return nil, entities.ErrorNotFound("Couldn't find this post")
	}

	if id == "1234" {
		return &Post{
			Id:    id,
			Title: "The Adventure time",
			Asset: nil,
		}, nil
	}

	assetId := "0195d2b8-bbbb-72a5-bcb1-27226d04b0c6"

	asset, err := client.GetAssetById(ctx, &proto.GetAssetRequest{Id: assetId})
	if err != nil {
		return nil, entities.ErrorNotFound(err)
	}

	return &Post{
		Id:    id,
		Title: "A quick brown fox jumps over the lazy dog",
		Asset: asset,
	}, nil
}
