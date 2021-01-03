package recipe

import (
	"errors"
)

type Validator struct {
	Recipes db
}

var ErrDisabled = errors.New("invalid recipe")

func (v Validator) Valid(id int) error {
	recipe, err := v.Recipes.Get(id)

	if err != nil {
		return err
	}

	if !recipe.Enabled {
		return ErrDisabled
	}

	return nil
}
