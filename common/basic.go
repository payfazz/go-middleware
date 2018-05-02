package common

import (
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/recovery"
)

// BasicPack return middleware pack that contain Logger, Recovery, and KV
func BasicPack() []middleware.Func {
	return BasicPackWithCb(nil, nil)
}

// BasicPackWithCb same as BasicPack, but you can provide logger and recovery callback
func BasicPackWithCb(
	loggerCb func(*logger.Event),
	recoveryCb func(*recovery.Event),
) []middleware.Func {
	return middleware.CompileList(
		logger.New(loggerCb),
		recovery.New(10, recoveryCb),
		kv.New(),
	)
}
