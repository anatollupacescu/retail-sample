package main

import (
	"sync/atomic"

	kitlog "github.com/go-kit/kit/log"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/types"
)

func newLoggerFactory(logger kitlog.Logger) types.LoggerFactory {
	counter := new(int32)

	return func() types.Logger {
		seq := atomic.AddInt32(counter, 1)
		return loggerWrapper{
			kitlog.With(logger, "request_id", seq),
		}
	}
}

type loggerWrapper struct {
	kitlog.Logger
}

func (lw loggerWrapper) Log(args ...interface{}) {
	_ = lw.Logger.Log(args...)
}
