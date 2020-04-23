package persistence

import (
	"context"
	"log"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/inventory"
)

type PgxInventoryStore struct {
	DB PgxDB
}

func (ps *PgxInventoryStore) Add(n inventory.Name) inventory.ID {
	var id int32
	err := ps.DB.QueryRow(context.Background(), "insert into inventory(name) values($1) returning id", n).Scan(&id)

	if err != nil {
		log.Print("inventory add ", err)
	}

	return inventory.ID(id)
}

func (ps *PgxInventoryStore) Find(n inventory.Name) (itemID inventory.ID) {
	rows, err := ps.DB.Query(context.Background(), "select id from inventory where name = $1", n)

	if err != nil {
		log.Print("inventory find", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			log.Print("inventory find scan", err)
		}

		itemID = inventory.ID(id)
	}

	return
}

func (ps *PgxInventoryStore) Get(id inventory.ID) inventory.Item {
	var name string
	err := ps.DB.QueryRow(context.Background(), "select name from inventory where id = $1", id).Scan(&name)

	if err != nil {
		log.Print("inventory get", err)
	}

	return inventory.Item{
		ID:   id,
		Name: inventory.Name(name),
	}
}

func (ps *PgxInventoryStore) List() (items []inventory.Item) {
	rows, err := ps.DB.Query(context.Background(), "select id, name from inventory")

	if err != nil {
		log.Print("inventory list", err)
	}

	defer rows.Close()

	for rows.Next() {
		var id int32
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Print("inventory list scan ", err)
			break
		}
		items = append(items, inventory.Item{
			ID:   inventory.ID(id),
			Name: inventory.Name(name),
		})
	}

	return
}
