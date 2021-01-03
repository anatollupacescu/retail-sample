package inventory

type Validator struct {
	Inventory db
}

func (v Validator) Validate(ids ...int) error {
	// check individual items
	for _, id := range ids {
		item, err := v.Inventory.Get(id)
		if err != nil {
			return err
		}

		if !item.Enabled {
			return ErrItemDisabled
		}
	}

	return nil
}
