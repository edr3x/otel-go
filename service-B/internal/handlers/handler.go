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

	return getPostByIDQuery(ctx, id)
}

type Post struct {
	ID    string                  `json:"id"`
	Title string                  `json:"title"`
	Asset *proto.GetAssetResponse `json:"asset"`
}

func getPostByIDQuery(ctx context.Context, id string) (*Post, error) {
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
			ID:    id,
			Title: "The Adventure time",
			Asset: nil,
		}, nil
	}

	asset, err := client.GetAssetById(ctx, &proto.GetAssetRequest{Id: id})
	if err != nil {
		return nil, entities.ErrorNotFound(err)
	}

	return &Post{
		ID:    id,
		Title: "A quick brown fox jumps over the lazy dog",
		Asset: asset,
	}, nil
}
