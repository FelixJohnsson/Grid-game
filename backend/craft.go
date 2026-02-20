package main

type CraftingRecipe struct {
	Name            string
	OutputItem      string
	OutputCount     int
	Ingredients     map[string]int
	RequiresShelter bool
}

var stoneAgeRecipes = map[string]CraftingRecipe{
	"Stone Knife": {
		Name:        "Stone Knife",
		OutputItem:  "Stone Knife",
		OutputCount: 1,
		Ingredients: map[string]int{
			stoneStockpileKey: 1,
			stickItemName:     1,
		},
		RequiresShelter: false,
	},
	"Stone Hammer": {
		Name:        "Stone Hammer",
		OutputItem:  "Stone Hammer",
		OutputCount: 1,
		Ingredients: map[string]int{
			stoneStockpileKey: 2,
			stickItemName:     1,
		},
		RequiresShelter: false,
	},
	"Stone Axe": {
		Name:        "Stone Axe",
		OutputItem:  "Stone Axe",
		OutputCount: 1,
		Ingredients: map[string]int{
			stoneStockpileKey: 2,
			stickItemName:     2,
		},
		RequiresShelter: false,
	},
	"Stone Spear": {
		Name:        "Stone Spear",
		OutputItem:  "Stone Spear",
		OutputCount: 1,
		Ingredients: map[string]int{
			stoneStockpileKey: 1,
			stickItemName:     2,
		},
		RequiresShelter: false,
	},
	"Stone Pickaxe": {
		Name:        "Stone Pickaxe",
		OutputItem:  "Stone Pickaxe",
		OutputCount: 1,
		Ingredients: map[string]int{
			stoneStockpileKey: 3,
			stickItemName:     2,
		},
		RequiresShelter: true,
	},
	"Food Box": {
		Name:        "Food Box",
		OutputItem:  "Food Box",
		OutputCount: 1,
		Ingredients: map[string]int{
			stickItemName:     4,
			stoneStockpileKey: 1,
		},
		RequiresShelter: true,
	},
	"Wooden Box": {
		Name:        "Wooden Box",
		OutputItem:  "Wooden Box",
		OutputCount: 1,
		Ingredients: map[string]int{
			stickItemName:     5,
			stoneStockpileKey: 1,
		},
		RequiresShelter: true,
	},
	"Wooden Crate": {
		Name:        "Wooden Crate",
		OutputItem:  "Wooden Crate",
		OutputCount: 1,
		Ingredients: map[string]int{
			stickItemName:     6,
			stoneStockpileKey: 2,
		},
		RequiresShelter: true,
	},
}

func copyIngredients(ingredients map[string]int) map[string]int {
	if len(ingredients) == 0 {
		return nil
	}

	copied := make(map[string]int, len(ingredients))
	for key, value := range ingredients {
		copied[key] = value
	}
	return copied
}

func GetCraftingRecipe(recipeName string) (CraftingRecipe, bool) {
	recipe, ok := stoneAgeRecipes[recipeName]
	if !ok {
		return CraftingRecipe{}, false
	}
	recipe.Ingredients = copyIngredients(recipe.Ingredients)
	return recipe, true
}

func (b *Brain) hasBaseOrOwnedItem(itemName string) bool {
	if itemName == "" {
		return false
	}
	if b.FindInOwnedItems(itemName) != nil {
		return true
	}
	if !b.Base.Claimed {
		return false
	}

	tile := b.Owner.WorldProvider.GetTile(b.Base.Location.X, b.Base.Location.Y)
	for _, item := range tile.Items {
		if item != nil && item.Name == itemName {
			return true
		}
	}
	return false
}

func (b *Brain) hasStockpileIngredients(ingredients map[string]int) bool {
	for key, needed := range ingredients {
		if needed <= 0 {
			continue
		}
		if b.Base.Stockpile[key] < needed {
			return false
		}
	}
	return true
}

func (b *Brain) consumeStockpileIngredients(ingredients map[string]int) {
	for key, needed := range ingredients {
		if needed <= 0 {
			continue
		}
		b.Base.Stockpile[key] -= needed
		if b.Base.Stockpile[key] < 0 {
			b.Base.Stockpile[key] = 0
		}
	}
}

func (b *Brain) refundStockpileIngredients(ingredients map[string]int) {
	for key, amount := range ingredients {
		if amount <= 0 {
			continue
		}
		b.Base.Stockpile[key] += amount
	}
}

// CraftItem crafts one configured recipe and drops result at base.
func (b *Brain) CraftItem(action TargetedAction) {
	recipeName := action.Target
	if recipeName == "" {
		switch action.Action {
		case CraftFoodBox:
			recipeName = "Food Box"
		case CraftStone:
			recipeName = "Stone Axe"
		default:
			return
		}
	}

	recipe, ok := GetCraftingRecipe(recipeName)
	if !ok {
		return
	}

	if !b.Base.Claimed {
		b.ClaimBaseTask(TargetedAction{Action: ClaimBase})
		return
	}
	if recipe.RequiresShelter && !b.PhysiologicalNeeds.HasShelter {
		return
	}

	b.ensureBaseStockpile()
	if !b.hasStockpileIngredients(recipe.Ingredients) {
		return
	}

	if b.hasBaseOrOwnedItem(recipe.OutputItem) {
		b.RemoveActionFromActionList(action)
		return
	}

	if b.Owner.Location != b.Base.Location {
		b.MotorCortexCurrentTask = MotorCortexAction{"Craft item", "Walk", b.Base.Location, false, false}
		return
	}

	outputCount := recipe.OutputCount
	if outputCount <= 0 {
		outputCount = 1
	}

	craftedItems := make([]*Item, 0, outputCount)
	for i := 0; i < outputCount; i++ {
		crafted := CreateNewItem(recipe.OutputItem)
		if crafted == nil {
			return
		}
		craftedItems = append(craftedItems, crafted)
	}

	b.consumeStockpileIngredients(recipe.Ingredients)

	added := 0
	for _, crafted := range craftedItems {
		if err := b.Owner.WorldProvider.AddItem(b.Base.Location.X, b.Base.Location.Y, crafted); err != nil {
			continue
		}
		added++
	}

	if added == 0 {
		b.refundStockpileIngredients(recipe.Ingredients)
		return
	}

	b.SyncKnownTileFromWorld(b.Base.Location)
	b.AddMemoryToLongTerm("Crafted item", recipe.OutputItem, b.Base.Location)
	b.RemoveActionFromActionList(action)
}

// GatherMaterials - Gather materials
func (b *Brain) GatherMaterials(action TargetedAction) {

}

// BuildShelter - Build a shelter
func (b *Brain) BuildShelter(action TargetedAction) {

}
