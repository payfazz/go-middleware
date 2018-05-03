package common

import (
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/recovery"
)

// BasicPack return middleware pack that contain Logger, Recovery, and KV
func BasicPack() []middleware.Func {
	return AdvancePack(nil, nil, 10)
}

// AdvancePack same as BasicPack, but you can provide logger callback, recovery callback, and stackTraceDepth
func AdvancePack(
	loggerCb func(*logger.Event),
	recoveryCb func(*recovery.Event),
	stackTraceDepth int,
) []middleware.Func {
	return middleware.CompileList(
		logger.New(loggerCb),
		recovery.New(stackTraceDepth, recoveryCb),
		kv.New(),
	)
}
