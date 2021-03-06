package persistence

import (
	"context"

	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/domain/retail/recipe"
)

type RecipePgxStore struct {
	DB pgx.Tx
}

func (pr *RecipePgxStore) Save(r recipe.DTO) error {
	sql := "update recipe set enabled=$1 where id=$2"

	tag, err := pr.DB.Exec(context.Background(), sql, r.Enabled, r.ID)

	if err != nil {
		return errors.Wrapf(ErrDB, "save recipe: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return recipe.ErrNotFound
	}

	return nil
}

func (pr *RecipePgxStore) Add(r recipe.DTO) (int, error) {
	sql := "insert into recipe(name) values($1) returning id"

	var recipeID int

	err := pr.DB.QueryRow(context.Background(), sql, r.Name).Scan(&recipeID)

	if err != nil {
		return 0, errors.Wrapf(ErrDB, "add recipe: %v", err)
	}

	sql = "insert into recipe_ingredient(recipeid, inventoryid, quantity) values($1, $2, $3)"

	for _, i := range r.Ingredients {
		_, err = pr.DB.Exec(context.Background(), sql, recipeID, i.ID, i.Qty)

		if err != nil {
			return 0, errors.Wrapf(ErrDB, "add recipe ingredient: %v", err)
		}
	}

	return recipeID, nil
}

func (pr *RecipePgxStore) Get(recipeID int) (r recipe.DTO, err error) {
	sql := "select name, enabled from recipe where id = $1"

	err = pr.DB.QueryRow(context.Background(), sql, recipeID).Scan(&r.Name, &r.Enabled)

	switch err {
	case nil:
	case pgx.ErrNoRows:
		return recipe.DTO{}, recipe.ErrNotFound
	default:
		return recipe.DTO{}, errors.Wrapf(ErrDB, "get recipe: %v", err)
	}

	sql = "select inventoryid, quantity from recipe_ingredient where recipeid = $1"

	rows, err := pr.DB.Query(context.Background(), sql, recipeID)

	if err != nil {
		return recipe.DTO{}, errors.Wrapf(ErrDB, "get recipe ingredients: %v", err)
	}

	defer rows.Close()

	var ingredients []recipe.InventoryItem

	for rows.Next() {
		var (
			itemid int64
			qty    int16
		)

		if err = rows.Scan(&itemid, &qty); err != nil {
			return recipe.DTO{}, errors.Wrapf(ErrDB, "scan recipe ingredients: %v", err)
		}

		ingredients = append(ingredients, recipe.InventoryItem{
			ID:  int(itemid),
			Qty: int(qty),
		})
	}

	r.ID = recipeID
	r.Ingredients = ingredients

	return r, nil
}

func (pr *RecipePgxStore) Find(name string) (recipe.DTO, error) {
	sql := "select id, name, enabled from recipe where name = $1"

	var r recipe.DTO

	err := pr.DB.QueryRow(context.Background(), sql, name).Scan(&r.ID, &r.Name, &r.Enabled)

	switch err {
	case nil:
	case pgx.ErrNoRows:
		return r, recipe.ErrNotFound
	default:
		return r, errors.Wrapf(ErrDB, "get recipe: %v", err)
	}

	return r, nil
}

func (pr *RecipePgxStore) List() ([]recipe.Recipe, error) {
	sql := `SELECT
						r.id,
						r.name,
						ri.inventoryid,
						ri.quantity,
						r.enabled
					FROM
						recipe_ingredient ri,
						recipe r
					WHERE
						ri.recipeid = r.id`

	rows, err := pr.DB.Query(context.Background(), sql)

	if err != nil {
		return nil, errors.Wrapf(ErrDB, "list recipes: %v", err)
	}

	defer rows.Close()

	type key struct {
		id      int
		name    string
		enabled bool
	}

	recipeRecords := make(map[key][]recipe.InventoryItem)

	for rows.Next() {
		var (
			recipeID     int
			name         string
			ingredientID int
			qty          int
			enabled      bool
		)

		if err := rows.Scan(&recipeID, &name, &ingredientID, &qty, &enabled); err != nil {
			return nil, errors.Wrapf(ErrDB, "scan recipes: %v", err)
		}

		recipeKey := key{id: recipeID, name: name, enabled: enabled}
		i := recipe.InventoryItem{ID: ingredientID, Qty: qty}

		recipeRecords[recipeKey] = append(recipeRecords[recipeKey], i)
	}

	recipes := make([]recipe.Recipe, 0, len(recipeRecords))

	for k, v := range recipeRecords {
		recipes = append(recipes, recipe.Recipe{
			ID:          k.id,
			Name:        k.name,
			Enabled:     k.enabled,
			Ingredients: v,
			DB:          pr,
		})
	}

	return recipes, nil
}
