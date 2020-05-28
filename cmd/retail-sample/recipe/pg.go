package recipe

import (
	"context"

	"github.com/jackc/pgconn"
	pgx "github.com/jackc/pgx/v4"
	"github.com/pkg/errors"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

var DBErr = errors.New("postgres")

type PgxDB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
}

type PgxStore struct {
	DB PgxDB
}

func (pr *PgxStore) Save(r recipe.Recipe) error {
	sql := "update recipe set enabled=$1 where id=$2"

	tag, err := pr.DB.Exec(context.Background(), sql, r.Enabled, r.ID)

	if err != nil {
		return errors.Wrapf(DBErr, "save recipe: %v", err)
	}

	if tag.RowsAffected() != 1 {
		return recipe.ErrRecipeNotFound
	}

	return nil
}

func (pr *PgxStore) Add(r recipe.Recipe) (recipe.ID, error) {
	sql := "insert into recipe(name) values($1) returning id"

	var (
		recipeID int32
		zeroID   = recipe.ID(recipeID)
	)

	err := pr.DB.QueryRow(context.Background(), sql, r.Name).Scan(&recipeID)

	if err != nil {
		return zeroID, errors.Wrapf(DBErr, "add recipe: %v", err)
	}

	sql = "insert into recipe_ingredient(recipeid, inventoryid, quantity) values($1, $2, $3)"

	for _, i := range r.Ingredients {
		_, err = pr.DB.Exec(context.Background(), sql, recipeID, i.ID, i.Qty)

		if err != nil {
			return zeroID, errors.Wrapf(DBErr, "add recipe ingredient: %v", err)
		}
	}

	return recipe.ID(recipeID), nil
}

func (pr *PgxStore) Get(recipeID recipe.ID) (r recipe.Recipe, err error) {
	sql := "select name, enabled from recipe where id = $1"

	err = pr.DB.QueryRow(context.Background(), sql, recipeID).Scan(&r.Name, &r.Enabled)

	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		return r, recipe.ErrRecipeNotFound
	default:
		return r, errors.Wrapf(DBErr, "get recipe: %v", err)
	}

	sql = "select inventoryid, quantity from recipe_ingredient where recipeid = $1"

	rows, err := pr.DB.Query(context.Background(), sql, recipeID)

	if err != nil {
		return r, errors.Wrapf(DBErr, "get recipe ingredients: %v", err)
	}

	defer rows.Close()

	var ingredients []recipe.Ingredient

	for rows.Next() {
		var itemid int64
		var qty int16

		if err = rows.Scan(&itemid, &qty); err != nil {
			return r, errors.Wrapf(DBErr, "scan recipe ingredients: %v", err)
		}

		ingredients = append(ingredients, recipe.Ingredient{
			ID:  int(itemid),
			Qty: int(qty),
		})
	}

	r.ID = recipeID
	r.Ingredients = ingredients

	return
}

func (pr *PgxStore) List() (recipes []recipe.Recipe, err error) {
	sql := `SELECT
						r.id,
						r.name,
						i.id,
						ri.quantity
					FROM
						recipe_ingredient ri,
						recipe r,
						inventory i
					WHERE
						ri.recipeid = r.id
						AND ri.inventoryid = i.id`

	rows, err := pr.DB.Query(context.Background(), sql)

	if err != nil {
		return recipes, errors.Wrapf(DBErr, "list recipes: %v", err)
	}

	defer rows.Close()

	type key struct {
		id   int64
		name string
	}

	type ingredients []recipe.Ingredient

	recipeRecords := make(map[key]ingredients)

	for rows.Next() {
		var (
			recipeID     int64
			name         string
			ingredientID int64
			qty          int16
		)

		if err := rows.Scan(&recipeID, &name, &ingredientID, &qty); err != nil {
			return recipes, errors.Wrapf(DBErr, "scan recipes: %v", err)
		}

		key := key{id: recipeID, name: name}
		i := recipe.Ingredient{ID: int(ingredientID), Qty: int(qty)}

		recipeRecords[key] = append(recipeRecords[key], i)
	}

	for k, v := range recipeRecords {
		recipes = append(recipes, recipe.Recipe{
			ID:          recipe.ID(k.id),
			Name:        recipe.Name(k.name),
			Ingredients: v,
		})
	}

	return
}
