// Package common provide pack that commonly used in every middleware chain
package common

import (
	middleware "github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/paniclogger"
)

// BasicPack return middleware pack that contain Logger, PanicLogger, and KV
func BasicPack() []middleware.Func {
	return AdvancePack(nil, nil, 10)
}

// AdvancePack same as BasicPack, but you can provide logger callback, paniclogger callback, and stackTraceDepth
func AdvancePack(
	loggerCb func(logger.Event),
	pannicloggerCb func(paniclogger.Event),
	stackTraceDepth int,
) []middleware.Func {
	return middleware.CompileList(
		logger.New(loggerCb),
		paniclogger.New(stackTraceDepth, pannicloggerCb),
		kv.New(),
	)
}
