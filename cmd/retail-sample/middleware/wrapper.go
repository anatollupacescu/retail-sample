package middleware

import (
	"fmt"
)

type Wrapper struct {
	LoggerFactory              LoggerFactory
	PersistenceProviderFactory PersistenceProviderFactory
}

type appHandler func(PersistenceProvider) error

func (ia Wrapper) Exec(action string, f appHandler) (err error) {
	logger := ia.LoggerFactory()

	logger.Log("msg", "enter", "action", action)

	provider := ia.PersistenceProviderFactory.New()

	wrapped := recoverHandler(middleware(f, provider))

	err = wrapped()

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

func middleware(f appHandler, p PersistenceProvider) func() error {
	return func() error {
		return f(p)
	}
}

func recoverHandler(next func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()

		return next()
	}
}
