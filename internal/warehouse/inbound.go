package warehouse

import "errors"

type (
	InboundType string

	Stock struct {
		inventory             map[InboundType]int
		outboundConfiguration map[OutboundType]OutboundItem
	}

	InboundItem struct {
		Type InboundType
		Qty  int
	}
)

var ErrItemTypeNotFound = errors.New("type not found")

func (s Stock) Provision(item InboundItem) (int, error) {

	if !s.hasType(item.Type) {
		return 0, ErrItemTypeNotFound
	}

	currentQty := s.inventory[item.Type]
	currentQty += item.Qty

	s.inventory[item.Type] = currentQty

	return s.inventory[item.Type], nil
}

var ErrItemTypePresent = errors.New("item type present")

func (s *Stock) ConfigureInboundType(typeName string) error {
	typeToAdd := InboundType(typeName)

	if s.hasType(typeToAdd) {
		return ErrItemTypePresent
	}

	if s.inventory == nil {
		s.inventory = make(map[InboundType]int)
	}

	s.inventory[typeToAdd] = 0

	return nil
}

func (s Stock) hasType(itemType InboundType) bool {
	_, found := s.inventory[itemType]
	return found
}

func (s Stock) Quantity(typeName string) (int, error) {
	itemType := InboundType(typeName)
	qty, found := s.inventory[itemType]

	if !found {
		return 0, ErrItemTypePresent
	}

	return qty, nil
}

func (s Stock) ItemTypes() (r []string) {
	for key := range s.inventory {
		r = append(r, string(key))
	}
	return
}
