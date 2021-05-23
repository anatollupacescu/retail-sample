package inventory

type Validator struct {
	Inventory db
}

func (v Validator) Valid(id int) (bool, error) {
	item, err := v.Inventory.Get(id)

	if err != nil {
		return false, err
	}

	return item.Enabled, nil
}
