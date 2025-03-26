package handlers

import (
	"context"
	"time"

	"github.com/edr3x/otel-go/interfaces"
	"github.com/edr3x/otel-go/internal/entities"
	"github.com/edr3x/otel-go/otelx"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
)

type Handler struct {
	res interfaces.Responders
}

func NewHandler(res interfaces.Responders) *Handler {
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

	return h.res.JSON(c, struct {
		Id    string `json:"id"`
		Title string `json:"title"`
	}{
		Id:    id,
		Title: name,
	})
}

func getPost(ctx context.Context, id string) (string, error) {
	ctx, span := otelx.StartSpan(ctx)
	defer span.End()

	time.Sleep(2 * time.Second)

	return getPostByIdQuery(ctx, id)
}

func getPostByIdQuery(ctx context.Context, id string) (string, error) {
	_, span := otelx.StartSpan(ctx)
	defer span.End()

	if id == "1234" {
		return "The Adventure time", nil
	}

	if id == "666" {
		span.SetAttributes(attribute.String("id", "666"), attribute.String("extrainfo", "cannot find book of the beast"))
		return "", entities.ErrorNotFound("Couldn't find this post")
	}

	return "A quick brown fox jumps over the lazy dog", nil
}
