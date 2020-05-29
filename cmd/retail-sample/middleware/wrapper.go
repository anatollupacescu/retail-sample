package middleware

type Wrapper struct {
	LoggerFactory              LoggerFactory
	PersistenceProviderFactory PersistenceProviderFactory
}

func (ia Wrapper) Exec(action string, f func(PersistenceProvider) error) {
	logger := ia.LoggerFactory()

	logger.Log("msg", "enter", "action", action)
	defer logger.Log("msg", "exit")

	provider := ia.PersistenceProviderFactory.New()

	err := f(provider)

	if err != nil {
		logger.Log("error", err)
		logger.Log("msg", "rollback")

		ia.PersistenceProviderFactory.Rollback(provider)
		return
	}

	logger.Log("msg", "commit")
	ia.PersistenceProviderFactory.Commit(provider)
}
