package middleware

import (
	"sync/atomic"

	kitlog "github.com/go-kit/kit/log"
)

type (
	Logger interface {
		Log(keyvals ...interface{})
	}

	NewLoggerFunc func() Logger
)

func BuildNewLoggerFunc(logger kitlog.Logger) NewLoggerFunc {
	counter := new(int32)

	return func() Logger {
		seq := atomic.AddInt32(counter, 1)

		return loggerWrapper{
			kitlog.With(logger, "request_id", seq),
		}
	}
}

func WrapLogger(logger kitlog.Logger) Logger {
	return loggerWrapper{
		Logger: kitlog.With(logger),
	}
}

type loggerWrapper struct {
	kitlog.Logger
}

func (lw loggerWrapper) Log(args ...interface{}) {
	_ = lw.Logger.Log(args...)
}
