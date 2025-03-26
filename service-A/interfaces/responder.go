package interfaces

import "github.com/labstack/echo/v4"

type Responders interface {
	JSON(c echo.Context, i any, code ...int) error
}
