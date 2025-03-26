package handlers

import (
	"context"
	"time"

	"github.com/edr3x/otel-go/pkg/entities"
	"github.com/edr3x/otel-go/pkg/entities/responders"
	"github.com/edr3x/otel-go/pkg/otelx"
	"github.com/edr3x/otel-go/service-a/internal/services"
	"github.com/labstack/echo/v4"
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

	name, err := getName(ctx, id)
	if err != nil {
		return err
	}

	post, err := services.GetUsersPosts(ctx, "666as")
	if err != nil {
		return entities.ErrorNotFound(err)
	}

	time.Sleep(1 * time.Second)

	return h.res.JSON(c, struct {
		Id   string                  `json:"id"`
		Name string                  `json:"name"`
		Post services.GetPostPayload `json:"post"`
	}{
		Id:   id,
		Name: name,
		Post: *post,
	})
}

func getName(ctx context.Context, id string) (string, error) {
	ctx, span := otelx.StartSpan(ctx)
	defer span.End()

	time.Sleep(2 * time.Second)

	return getNameByIdQuery(ctx, id)
}

func getNameByIdQuery(ctx context.Context, id string) (string, error) {
	_, span := otelx.StartSpan(ctx)
	defer span.End()

	time.Sleep(3 * time.Second)
	if id == "1234" {
		return "John Doe", nil
	}

	if id == "222" {
		return "", entities.ErrorBadRequest("cannot provide this id")
	}

	return "Hylos", nil
}
