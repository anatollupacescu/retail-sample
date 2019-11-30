package warehouse

import "errors"

type (
	OutboundItem struct {
		ItemType ItemType
		Qty      int
	}

	FinishedProduct struct {
		name  string
		Items []OutboundItem
	}
)

var (
	ErrNameNotProvided        = errors.New("name not provided")
	ErrItemsNotProvided       = errors.New("items not provided")
	ErrZeroQuantityNotAllowed = errors.New("zero quantity not allowed")
)

func (s *Stock) AddFinishedProduct(name string, items []OutboundItem) error {
	if len(name) == 0 {
		return ErrNameNotProvided
	}

	if len(items) == 0 {
		return ErrItemsNotProvided
	}

	for _, item := range items {
		if !s.hasType(item.ItemType) {
			return ErrItemTypeNotFound
		}
		if item.Qty == 0 {
			return ErrZeroQuantityNotAllowed
		}
	}

	newFinishedProduct := FinishedProduct{
		name:  name,
		Items: items,
	}
	s.finishedProducts = append(s.finishedProducts, newFinishedProduct)

	return nil
}

func (s *Stock) FinishedProducts() []FinishedProduct {
	return s.finishedProducts
}
