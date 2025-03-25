package handlers

import (
	"context"
	"time"

	"github.com/edr3x/otel-go/interfaces"
	"github.com/edr3x/otel-go/internal/entities"
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

	time.Sleep(1 * time.Second)

	return h.res.JSON(c, struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{
		Id:   id,
		Name: name,
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
		return "", entities.ErrorBadRequest("cannot provide this id", map[string]any{
			"id":      id,
			"context": "you can never provide 22222",
		})
	}

	return "Hylos", nil
}
