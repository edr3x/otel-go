package handlers

import (
	"context"
	"time"

	"github.com/edr3x/otel-go/interfaces"
	"github.com/edr3x/otel-go/internal/entities"
	"github.com/edr3x/otel-go/internal/services"
	"github.com/edr3x/otel-go/otelx"
	"github.com/labstack/echo/v4"
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

	name, err := getName(ctx, id)
	if err != nil {
		return err
	}

	post, err := services.GetUsersPosts(ctx, "666")
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
