package recipe

type Validator struct {
	Recipes db
}

func (v Validator) Valid(id int) (bool, error) {
	recipe, err := v.Recipes.Get(id)

	if err != nil {
		return false, err
	}

	return recipe.Enabled, nil
}
