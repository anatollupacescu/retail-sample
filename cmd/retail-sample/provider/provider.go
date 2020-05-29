package provider

import "github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"

func NewPersistenceFactory(dbConn string, inMemory bool) middleware.PersistenceProviderFactory {
	if inMemory {
		return newInMemoryPersistentFactory()
	}

	return newPersistenceFactory(dbConn)
}
