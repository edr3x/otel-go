package responders

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Responders interface {
	JSON(c echo.Context, i any, code ...int) error
}

type res struct {
	logger *zap.Logger
}

func NewResponder() Responders {
	return &res{
		logger: zap.L(),
	}
}
