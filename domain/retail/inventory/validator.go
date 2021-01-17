package inventory

import "errors"

type Validator struct {
	Inventory db
}

var ErrItemDisabled = errors.New("item is disabled")

func (v Validator) Validate(id int) error {
	item, err := v.Inventory.Get(id)

	if err != nil {
		return err
	}

	if !item.Enabled {
		return ErrItemDisabled
	}

	return nil
}
