package warehouse

type Stock struct {
	inventory    Inventory
	provisionLog ProvisionLog
	orderLog     OrderLog
	recipeBook   RecipeBook
	data         map[int]int
}

func NewStock(inv Inventory, recipeBook RecipeBook, provisionLog ProvisionLog, orderLog OrderLog) Stock {
	return Stock{
		provisionLog: provisionLog,
		orderLog:     orderLog,
		inventory:    inv,
		recipeBook:   recipeBook,
		data:         make(map[int]int),
	}
}

func NewStockWithData(inv Inventory, recipeBook RecipeBook, provisionLog ProvisionLog, orderLog OrderLog, d map[int]int) Stock {
	return Stock{
		provisionLog: provisionLog,
		orderLog:     orderLog,
		inventory:    inv,
		recipeBook:   recipeBook,
		data:         d,
	}
}
