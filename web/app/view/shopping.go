package view

import (
	"fmt"
	"log"
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

	templateBytes := appTemplate.MustAsset("web/template/shopping_list.html")
	t := template.New("shopping_list.html")
	t.Parse(string(templateBytes[:]))
	t.Execute(w, shoppingList)
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

	// Add recipe ingredients to shopping list
	for _, ingredient := range recipe.Ingredients {
		item := shopping.ShoppingItem{}
		item.FromIngredient(ingredient, recipeId, recipe.Name)
		shoppingList.Items = append(shoppingList.Items, item)
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