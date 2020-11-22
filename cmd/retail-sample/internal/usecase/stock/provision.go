package stock

type UpdateDTO struct {
	InventoryItemID int
	Qty             int
}

type Position struct {
	ID   int
	Name string
	Qty  int
}

func (o *Stock) Provision(dto UpdateDTO) (Position, error) {
	o.logger.Info("provision", "enter")

	_, err := o.stock.Provision(dto.InventoryItemID, dto.Qty)
	if err != nil {
		o.logger.Error("provision", "call domain layer", err)
		return Position{}, err
	}

	qty, err := o.stockDB.Quantity(dto.InventoryItemID)

	if err != nil {
		o.logger.Error("provision", "call domain layer to retrieve quantity", err)
		return Position{}, err
	}

	item, err := o.inventoryDB.Get(dto.InventoryItemID)
	if err != nil {
		o.logger.Error("provision", "call domain layer to retrieve stock position", err)
		return Position{}, err
	}

	o.logger.Error("provision", "success", err)

	pos := Position{
		ID:   item.ID,
		Name: item.Name,
		Qty:  qty,
	}

	return pos, nil
}
