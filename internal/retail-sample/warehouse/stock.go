package warehouse

type Stock struct {
	inventory    Inventory
	provisionLog ProvisionLog
	orders       Orders
	recipeBook   RecipeBook
	data         map[int]int
}

func NewStock(inventory Inventory, recipeBook RecipeBook, provisionLog ProvisionLog, orders Orders) Stock {
	return Stock{
		provisionLog: provisionLog,
		orders:       orders,
		inventory:    inventory,
		recipeBook:   recipeBook,
		data:         make(map[int]int),
	}
}

func NewStockWithData(inv Inventory, recipeBook RecipeBook, provisionLog ProvisionLog, orders Orders, d map[int]int) Stock {
	return Stock{
		provisionLog: provisionLog,
		orders:       orders,
		inventory:    inv,
		recipeBook:   recipeBook,
		data:         d,
	}
}
