package olog

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	Logger = zap.NewExample()
}
