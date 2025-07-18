package view

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"github.com/mryachanin/cookbook/api/recipe"
	"github.com/mryachanin/cookbook/api/shopping"
	appTemplate "github.com/mryachanin/cookbook/web/template"
	"github.com/julienschmidt/httprouter"
	"github.com/rhinoman/couchdb-go"
	"html/template"
	"net/http"
	"time"
)

const shoppingListId = "shopping-list"
const exclusionListId = "exclusion-list"

func GetShoppingList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
	shoppingList := shopping.ShoppingList{}
	_, err := db.Read(shoppingListId, &shoppingList, nil)
	if err != nil {
		// If shopping list doesn't exist, create empty one
		shoppingList = shopping.ShoppingList{
			Id:        shoppingListId,
			Items:     []shopping.ShoppingItem{},
			UpdatedAt: time.Now(),
		}
	}

	// Get exclusion list and filter out excluded ingredients
	exclusionList := shopping.ExclusionList{}
	_, err = db.Read(exclusionListId, &exclusionList, nil)
	if err == nil {
		// Filter items that are not excluded
		filteredItems := []shopping.ShoppingItem{}
		for _, item := range shoppingList.Items {
			if !isIngredientExcluded(item.Name, exclusionList.ExcludedIngredients) {
				filteredItems = append(filteredItems, item)
			}
		}
		shoppingList.Items = filteredItems
	}

	templateBytes := appTemplate.MustAsset("web/template/shopping_list.html")
	t := template.New("shopping_list.html")
	t.Parse(string(templateBytes[:]))
	t.Execute(w, shoppingList)
}

func isIngredientExcluded(ingredientName string, excludedIngredients []string) bool {
	for _, excluded := range excludedIngredients {
		if strings.EqualFold(ingredientName, excluded) {
			return true
		}
	}
	return false
}

func AddToShoppingList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
	recipeId := ps.ByName("id")
	
	// Get the recipe
	recipe := recipe.Recipe{}
	_, err := db.Read(recipeId, &recipe, nil)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	// Get existing shopping list
	shoppingList := shopping.ShoppingList{}
	rev, err := db.Read(shoppingListId, &shoppingList, nil)
	if err != nil {
		// Create new shopping list if it doesn't exist
		shoppingList = shopping.ShoppingList{
			Id:        shoppingListId,
			Items:     []shopping.ShoppingItem{},
			UpdatedAt: time.Now(),
		}
		rev = ""
	}

	// Get exclusion list to filter out excluded ingredients
	exclusionList := shopping.ExclusionList{}
	_, err = db.Read(exclusionListId, &exclusionList, nil)
	
	// Add recipe ingredients to shopping list (excluding any that are in the exclusion list)
	for _, ingredient := range recipe.Ingredients {
		if err != nil || !isIngredientExcluded(ingredient.Name, exclusionList.ExcludedIngredients) {
			item := shopping.ShoppingItem{}
			item.FromIngredient(ingredient, recipeId, recipe.Name)
			
			// Check if this ingredient already exists (match by name only)
			found := false
			for i, existingItem := range shoppingList.Items {
				if strings.EqualFold(existingItem.Name, item.Name) {
					// Combine quantities
					shoppingList.Items[i].Quantity = combineQuantities(existingItem.Quantity, item.Quantity)
					// Keep the original prep and note, add recipe name if different
					if !strings.Contains(existingItem.RecipeName, recipe.Name) {
						if existingItem.RecipeName != "" {
							shoppingList.Items[i].RecipeName = existingItem.RecipeName + ", " + recipe.Name
						} else {
							shoppingList.Items[i].RecipeName = recipe.Name
						}
					}
					found = true
					break
				}
			}
			
			if !found {
				shoppingList.Items = append(shoppingList.Items, item)
			}
		}
	}
	
	shoppingList.UpdatedAt = time.Now()

	// Save shopping list
	_, err = db.Save(&shoppingList, shoppingListId, rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ingredients added to shopping list"))
}

func RemoveFromShoppingList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB max memory
	if err != nil {
		// If not multipart, try regular form parsing
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	itemIndex := r.FormValue("index")
	if itemIndex == "" {
		http.Error(w, "Item index required", http.StatusBadRequest)
		return
	}

	// Get shopping list
	shoppingList := shopping.ShoppingList{}
	rev, err := db.Read(shoppingListId, &shoppingList, nil)
	if err != nil {
		http.Error(w, "Shopping list not found", http.StatusNotFound)
		return
	}

	// Parse index and remove item
	var index int
	if _, err := fmt.Sscanf(itemIndex, "%d", &index); err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	if index < 0 || index >= len(shoppingList.Items) {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}

	// Remove item from slice
	shoppingList.Items = append(shoppingList.Items[:index], shoppingList.Items[index+1:]...)
	shoppingList.UpdatedAt = time.Now()

	// Save updated shopping list
	_, err = db.Save(&shoppingList, shoppingListId, rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item removed from shopping list"))
}

func ToggleShoppingItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB max memory
	if err != nil {
		// If not multipart, try regular form parsing
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	itemIndex := r.FormValue("index")
	log.Printf("Toggle shopping item - received index: %s", itemIndex)
	if itemIndex == "" {
		log.Printf("Toggle shopping item - no index provided")
		http.Error(w, "Item index required", http.StatusBadRequest)
		return
	}

	// Get shopping list
	shoppingList := shopping.ShoppingList{}
	rev, err := db.Read(shoppingListId, &shoppingList, nil)
	if err != nil {
		http.Error(w, "Shopping list not found", http.StatusNotFound)
		return
	}

	// Parse index and toggle item
	var index int
	if _, err := fmt.Sscanf(itemIndex, "%d", &index); err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	if index < 0 || index >= len(shoppingList.Items) {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}

	// Toggle checked status
	shoppingList.Items[index].Checked = !shoppingList.Items[index].Checked
	shoppingList.UpdatedAt = time.Now()

	// Save updated shopping list
	_, err = db.Save(&shoppingList, shoppingListId, rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Item toggled"))
}

func ExcludeCheckedItems(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
	// Get shopping list
	shoppingList := shopping.ShoppingList{}
	_, err := db.Read(shoppingListId, &shoppingList, nil)
	if err != nil {
		http.Error(w, "Shopping list not found", http.StatusNotFound)
		return
	}

	// Get or create exclusion list
	exclusionList := shopping.ExclusionList{}
	rev, err := db.Read(exclusionListId, &exclusionList, nil)
	if err != nil {
		// Create new exclusion list if it doesn't exist
		exclusionList = shopping.ExclusionList{
			Id:                  exclusionListId,
			ExcludedIngredients: []string{},
			UpdatedAt:           time.Now(),
		}
		rev = ""
	}

	// Collect checked ingredient names to exclude
	var checkedIngredients []string
	var uncheckedItems []shopping.ShoppingItem

	for _, item := range shoppingList.Items {
		if item.Checked {
			// Add to exclusion list (check for duplicates)
			found := false
			for _, excluded := range exclusionList.ExcludedIngredients {
				if strings.EqualFold(item.Name, excluded) {
					found = true
					break
				}
			}
			if !found {
				exclusionList.ExcludedIngredients = append(exclusionList.ExcludedIngredients, item.Name)
				checkedIngredients = append(checkedIngredients, item.Name)
			}
		} else {
			// Keep unchecked items
			uncheckedItems = append(uncheckedItems, item)
		}
	}

	if len(checkedIngredients) == 0 {
		http.Error(w, "No checked items to exclude", http.StatusBadRequest)
		return
	}

	// Save updated exclusion list
	exclusionList.UpdatedAt = time.Now()
	_, err = db.Save(&exclusionList, exclusionListId, rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update shopping list to remove checked items
	shoppingList.Items = uncheckedItems
	shoppingList.UpdatedAt = time.Now()
	shoppingRev, err := db.Read(shoppingListId, &shopping.ShoppingList{}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Save(&shoppingList, shoppingListId, shoppingRev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Excluded %d ingredients", len(checkedIngredients))))
}

func GetExcludedIngredients(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
	exclusionList := shopping.ExclusionList{}
	_, err := db.Read(exclusionListId, &exclusionList, nil)
	if err != nil {
		// If exclusion list doesn't exist, create empty one
		exclusionList = shopping.ExclusionList{
			Id:                  exclusionListId,
			ExcludedIngredients: []string{},
			UpdatedAt:           time.Now(),
		}
	}

	templateBytes := appTemplate.MustAsset("web/template/excluded_ingredients.html")
	t := template.New("excluded_ingredients.html")
	t.Parse(string(templateBytes[:]))
	t.Execute(w, exclusionList)
}

func RemoveExcludedIngredient(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB max memory
	if err != nil {
		// If not multipart, try regular form parsing
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	ingredientName := r.FormValue("ingredient")
	if ingredientName == "" {
		http.Error(w, "Ingredient name required", http.StatusBadRequest)
		return
	}

	// Get exclusion list
	exclusionList := shopping.ExclusionList{}
	rev, err := db.Read(exclusionListId, &exclusionList, nil)
	if err != nil {
		http.Error(w, "Exclusion list not found", http.StatusNotFound)
		return
	}

	// Remove ingredient from exclusion list
	found := false
	filteredIngredients := []string{}
	for _, excluded := range exclusionList.ExcludedIngredients {
		if !strings.EqualFold(excluded, ingredientName) {
			filteredIngredients = append(filteredIngredients, excluded)
		} else {
			found = true
		}
	}

	if !found {
		http.Error(w, "Ingredient not found in exclusion list", http.StatusNotFound)
		return
	}

	// Update exclusion list
	exclusionList.ExcludedIngredients = filteredIngredients
	exclusionList.UpdatedAt = time.Now()

	// Save updated exclusion list
	_, err = db.Save(&exclusionList, exclusionListId, rev)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ingredient removed from exclusion list"))
}

func combineQuantities(qty1, qty2 string) string {
	if qty1 == "" {
		return qty2
	}
	if qty2 == "" {
		return qty1
	}

	// Try to extract numbers and units from both quantities
	num1, unit1 := parseQuantity(qty1)
	num2, unit2 := parseQuantity(qty2)

	// If units match or one is empty, we can add the numbers
	if unit1 == unit2 || unit1 == "" || unit2 == "" {
		combinedNum := num1 + num2
		unit := unit1
		if unit == "" {
			unit = unit2
		}
		
		// Format the combined quantity
		if combinedNum == float64(int(combinedNum)) {
			return fmt.Sprintf("%.0f %s", combinedNum, unit)
		} else {
			return fmt.Sprintf("%.1f %s", combinedNum, unit)
		}
	}

	// If units don't match, just concatenate with "+"
	return qty1 + " + " + qty2
}

func parseQuantity(quantity string) (float64, string) {
	// Regular expression to match number at the beginning
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?(?:/\d+)?)\s*(.*)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(quantity))
	
	if len(matches) < 3 {
		return 0, quantity // Return 0 and full string if no number found
	}

	numStr := matches[1]
	unit := strings.TrimSpace(matches[2])

	// Handle fractions like "1/2"
	if strings.Contains(numStr, "/") {
		parts := strings.Split(numStr, "/")
		if len(parts) == 2 {
			numerator, err1 := strconv.ParseFloat(parts[0], 64)
			denominator, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 == nil && err2 == nil && denominator != 0 {
				return numerator / denominator, unit
			}
		}
	}

	// Parse regular number
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, quantity
	}

	return num, unit
}