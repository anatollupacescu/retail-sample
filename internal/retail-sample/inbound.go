package warehouse

import "errors"

type (
	ItemType string

	Stock struct {
		stockItems       map[ItemType]int
		finishedProducts []FinishedProduct
	}

	InboundItem struct {
		Type ItemType
		Qty  int
	}
)

var ErrItemTypeNotFound = errors.New("type not found")

func (s Stock) Add(item InboundItem) (int, error) {
	if !s.hasType(item.Type) {
		return 0, ErrItemTypeNotFound
	}

	currentQty := s.stockItems[item.Type]
	currentQty += item.Qty

	s.stockItems[item.Type] = currentQty

	return s.stockItems[item.Type], nil
}

var ErrItemTypePresent = errors.New("item type present")

func (s *Stock) AddType(typeName string) error {
	typeToAdd := ItemType(typeName)

	if s.hasType(typeToAdd) {
		return ErrItemTypePresent
	}

	if s.stockItems == nil {
		s.stockItems = make(map[ItemType]int)
	}

	s.stockItems[typeToAdd] = 0

	return nil
}

func (s Stock) hasType(itemType ItemType) bool {
	_, found := s.stockItems[itemType]
	return found
}
