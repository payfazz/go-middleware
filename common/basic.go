package common

import (
	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/recovery"
)

// BasicPack return middleware pack that contain Logger, Recovery, and KV
func BasicPack() []middleware.Func {
	return middleware.CompileList(
		logger.New(nil),
		recovery.New(10, nil),
		kv.New(),
	)
}
