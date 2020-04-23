package persistence

import (
	"context"
	"errors"
	"log"

	"github.com/anatollupacescu/retail-sample/internal/retail-domain/recipe"
)

type PgxRecipeStore struct {
	DB PgxDB
}

var DbErr = errors.New("request to the database failedd")

func (pr *PgxRecipeStore) Add(r recipe.Recipe) (recipe.ID, error) {
	var (
		recipeID int32
		zeroID   = recipe.ID(recipeID)
	)

	sql := "insert into recipe(name) values($1) returning id"
	err := pr.DB.QueryRow(context.Background(), sql, r.Name).Scan(&recipeID)

	if err != nil {
		log.Print("recipe add", err)
		return zeroID, DbErr
	}

	sql = "insert into recipe_ingredient(recipeid, inventoryid, quantity) values($1, $2, $3)"
	for _, i := range r.Ingredients {
		_, err = pr.DB.Exec(context.Background(), sql, recipeID, i.ID, i.Qty)

		if err != nil {
			log.Print("recipe ingredient add", err)
			return zeroID, nil
		}
	}

	return recipe.ID(recipeID), nil
}

func (pr *PgxRecipeStore) Get(recipeID recipe.ID) (r recipe.Recipe) {
	sql := "select name from recipe where id = $1"

	var name string
	err := pr.DB.QueryRow(context.Background(), sql, recipeID).Scan(&name)

	if err != nil {
		log.Print("recipe get", err)
		return
	}

	sql = "select inventoryid, quantity from recipe_ingredient where recipeid = $1"

	rows, err := pr.DB.Query(context.Background(), sql, recipeID)

	if err != nil {
		log.Print("recipe get ingredients ", err)
		return
	}

	defer rows.Close()

	var ingredients []recipe.Ingredient

	for rows.Next() {
		var itemid int64
		var qty int16

		if err := rows.Scan(&itemid, &qty); err != nil {
			log.Print("recipe scan ingredients", err)
			return
		}

		ingredients = append(ingredients, recipe.Ingredient{
			ID:  int(itemid),
			Qty: int(qty),
		})
	}

	r.Ingredients = ingredients
	return
}

func (pr *PgxRecipeStore) List() (recipes []recipe.Recipe) {
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
		log.Print("recipe list", err)
		return nil
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
			log.Print("recipe list scan", err)
			return nil
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
