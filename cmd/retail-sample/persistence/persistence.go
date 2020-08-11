package persistence

import (
	"strings"

	"github.com/anatollupacescu/retail-sample/cmd/retail-sample/middleware"
)

func NewPersistenceFactory(dbConn string) middleware.PersistenceProviderFactory {
	if strings.TrimSpace(dbConn) == "" {
		return newInMemoryPersistentFactory()
	}

	return newPersistenceFactory(dbConn)
}
