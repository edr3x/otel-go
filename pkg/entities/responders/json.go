package responders

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type successResponse struct {
	Success bool `json:"success"`
	Payload any  `json:"payload"`
}

func (r *res) JSON(c echo.Context, i any, code ...int) error {
	status := 200
	span := trace.SpanFromContext(c.Request().Context())
	span.SetStatus(codes.Ok, "success")
	if len(code) > 0 {
		status = code[0]
	}
	return c.JSON(status, successResponse{Success: true, Payload: i})
}
