package middleware

import (
	"fmt"
)

type Middleware struct {
	NewLogger                  NewLoggerFunc
	PersistenceProviderFactory PersistenceProviderFactory
}

type appHandler func(PersistenceProvider) error

func (ia Middleware) Exec(action string, handler appHandler) (err error) {
	logger := ia.NewLogger()
	logger.Log("msg", "enter", "action", action)

	provider := ia.PersistenceProviderFactory.New()

	handlerCall := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()

		return handler(provider)
	}

	err = handlerCall()

	if err != nil {
		logger.Log("error", err)
		logger.Log("msg", "rollback")
		ia.PersistenceProviderFactory.Rollback(provider)

		return
	}

	logger.Log("msg", "commit")
	ia.PersistenceProviderFactory.Commit(provider)
	logger.Log("msg", "exit")

	return
}
