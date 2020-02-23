package warehouse

type Stock struct {
	inventory   Inventory
	inboundLog  InboundLog
	outboundLog OutboundLog
	recipeBook  RecipeBook
	data        map[int]int
}

func NewStock(log InboundLog, inv Inventory, recipeBook RecipeBook, outboundItemLog OutboundLog) Stock {
	return Stock{
		inboundLog:  log,
		outboundLog: outboundItemLog,
		inventory:   inv,
		recipeBook:  recipeBook,
	}
}

func NewStockWithData(log InboundLog, inv Inventory, recipeBook RecipeBook, outboundItemLog OutboundLog, d map[int]int) Stock {
	return Stock{
		inboundLog:  log,
		outboundLog: outboundItemLog,
		inventory:   inv,
		recipeBook:  recipeBook,
		data:        d,
	}
}
