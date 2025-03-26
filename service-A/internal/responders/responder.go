package responders

import (
	"github.com/edr3x/otel-go/interfaces"
	"go.uber.org/zap"
)

type res struct {
	logger *zap.Logger
}

func NewResponder() interfaces.Responders {
	return &res{
		logger: zap.L(),
	}
}
