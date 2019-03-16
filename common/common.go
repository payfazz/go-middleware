// Package common provide pack that commonly used in every middleware chain
package common

import (
	"log"
	"net/http"

	"github.com/payfazz/go-middleware"
	"github.com/payfazz/go-middleware/common/kv"
	"github.com/payfazz/go-middleware/common/logger"
	"github.com/payfazz/go-middleware/common/paniclogger"
)

// BasicPack return middleware pack that contain Logger, PanicLogger, and KV
func BasicPack() []func(http.HandlerFunc) http.HandlerFunc {
	return AdvancePack(nil, nil, 20)
}

// BasicPack2 do the same as BasicPack, but you can provide out and err logger.
func BasicPack2(out, err *log.Logger) []func(http.HandlerFunc) http.HandlerFunc {
	var loggerCb logger.Callback
	var panicLoggerCb paniclogger.Callback
	if out != nil {
		loggerCb = logger.DefaultLogger(out)
	}
	if err != nil {
		panicLoggerCb = paniclogger.DefaultLogger(err)
	}
	return AdvancePack(loggerCb, panicLoggerCb, 20)
}

// AdvancePack same as BasicPack, but you can provide logger callback, paniclogger callback, and stackTraceDepth
func AdvancePack(loggerCb logger.Callback, pannicloggerCb paniclogger.Callback, stackTraceDepth int) []func(http.HandlerFunc) http.HandlerFunc {
	return middleware.CompileList(
		paniclogger.New(stackTraceDepth, pannicloggerCb),
		logger.New(loggerCb),
		kv.New(),
	)
}
