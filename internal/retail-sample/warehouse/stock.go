package warehouse

type Stock struct {
	Inventory   Inventory
	InboundLog  InboundLog
	OutboundLog OutboundLog
	RecipeBook  RecipeBook
	Data        map[int]int
}

func NewStock(log InboundLog, inv Inventory, recipeBook RecipeBook, outboundItemLog OutboundLog) Stock {
	return Stock{
		InboundLog:  log,
		OutboundLog: outboundItemLog,
		Inventory:   inv,
		RecipeBook:  recipeBook,
	}
}

func NewStockWithData(log InboundLog, inv Inventory, recipeBook RecipeBook, outboundItemLog OutboundLog, d map[int]int) Stock {
	return Stock{
		InboundLog:  log,
		OutboundLog: outboundItemLog,
		Inventory:   inv,
		RecipeBook:  recipeBook,
		Data:        d,
	}
}
