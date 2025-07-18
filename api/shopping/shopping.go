package shopping

import (
	"time"
	"github.com/mryachanin/cookbook/api/recipe"
)

type ShoppingItem struct {
	Id       string `json:"_id,omitempty"`
	Rev      string `json:"_rev,omitempty"`
	Quantity string `json:"quantity"`
	Name     string `json:"name"`
	Prep     string `json:"prep,omitempty"`
	Note     string `json:"note,omitempty"`
	RecipeId string `json:"recipe_id,omitempty"`
	RecipeName string `json:"recipe_name,omitempty"`
	AddedAt  time.Time `json:"added_at"`
	Checked  bool   `json:"checked"`
}

type ShoppingList struct {
	Id    string         `json:"_id,omitempty"`
	Rev   string         `json:"_rev,omitempty"`
	Items []ShoppingItem `json:"items"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// Convert recipe ingredient to shopping item
func (item *ShoppingItem) FromIngredient(ingredient recipe.Ingredient, recipeId, recipeName string) {
	item.Quantity = ingredient.Quantity
	item.Name = ingredient.Name
	item.Prep = ingredient.Prep
	item.Note = ingredient.Note
	item.RecipeId = recipeId
	item.RecipeName = recipeName
	item.AddedAt = time.Now()
	item.Checked = false
}