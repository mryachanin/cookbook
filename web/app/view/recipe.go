package view

import (
  "github.com/mryachanin/cookbook/api/recipe"
  appTemplate "github.com/mryachanin/cookbook/web/template"
  "github.com/julienschmidt/httprouter"
  "github.com/rhinoman/couchdb-go"
  "github.com/google/uuid"
  "errors"
  "html/template"
  "net/http"
)

func GetRecipe(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
  recipeId := ps.ByName("id")

  if len(recipeId) == 0 {
    err := errors.New("Recipe ID must be specified")
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  recipe := recipe.Recipe{}
  db.Read(recipeId, &recipe, nil)

  templateBytes := appTemplate.MustAsset("web/template/view_recipe.html")
  t := template.New("view_recipe.html")
  t.Parse(string(templateBytes[:]))
  t.Execute(w, recipe)
}

func CreateRecipe(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
  templateBytes := appTemplate.MustAsset("web/template/create_recipe.html")
  t := template.New("create_recipe.html")
  t.Parse(string(templateBytes[:]))
  t.Execute(w, nil)
}

func PostRecipe(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
  err := r.ParseForm()
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // Create a new recipe
  newRecipe := recipe.Recipe{
    Name: r.FormValue("r_name"),
    Yields: r.FormValue("yields"),
    KeepsFor: r.FormValue("keeps_for"),
    PrepTime: r.FormValue("prep_time"),
    TotalTime: r.FormValue("total_time"),
  }

  // Process ingredients
  quantities := r.Form["i_quantity[]"]
  names := r.Form["i_name[]"]
  preps := r.Form["i_prep[]"]
  notes := r.Form["i_note[]"]

  for i := 0; i < len(quantities) && i < len(names); i++ {
    ingredient := recipe.Ingredient{
      Quantity: getValue(quantities, i),
      Name:     getValue(names, i),
      Prep:     getValue(preps, i),
      Note:     getValue(notes, i),
    }
    newRecipe.Ingredients = append(newRecipe.Ingredients, ingredient)
  }

  // Process instructions
  steps := r.Form["step[]"]
  for _, step := range steps {
    if step != "" {
      newRecipe.Instructions = append(newRecipe.Instructions, step)
    }
  }

  // Process notes
  recipeNotes := r.Form["note[]"]
  for _, note := range recipeNotes {
    if note != "" {
      newRecipe.Notes = append(newRecipe.Notes, note)
    }
  }

  // Process tags
  tags := r.Form["tag[]"]
  for _, tag := range tags {
    if tag != "" {
      newRecipe.Tags = append(newRecipe.Tags, tag)
    }
  }

  // Process source
  if r.FormValue("s_name") != "" {
    sourceType := recipe.Type(r.FormValue("s_type"))
    newRecipe.Source = recipe.Source{
      Name: r.FormValue("s_name"),
      Type: sourceType,
      Link: r.FormValue("s_link"),
    }
  }

  // Save to database
  _, err = db.Save(&newRecipe, uuid.New().String(), "")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Show success page
  templateBytes := appTemplate.MustAsset("web/template/create_recipe.html")
  t := template.New("create_recipe.html")
  t.Parse(string(templateBytes[:]))
  t.Execute(w, map[string]bool{"Success": true})
}

func EditRecipe(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
  recipeId := ps.ByName("id")

  if len(recipeId) == 0 {
    err := errors.New("Recipe ID must be specified")
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  recipe := recipe.Recipe{}
  _, err := db.Read(recipeId, &recipe, nil)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  templateBytes := appTemplate.MustAsset("web/template/edit_recipe.html")
  t := template.New("edit_recipe.html")
  
  // Add template function for adding numbers
  funcMap := template.FuncMap{
    "add": func(a, b int) int {
      return a + b
    },
  }
  t = t.Funcs(funcMap)
  
  t.Parse(string(templateBytes[:]))
  t.Execute(w, recipe)
}

func UpdateRecipe(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
  recipeId := ps.ByName("id")

  if len(recipeId) == 0 {
    err := errors.New("Recipe ID must be specified")
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Get existing recipe to preserve ID and Rev
  existingRecipe := recipe.Recipe{}
  rev, err := db.Read(recipeId, &existingRecipe, nil)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  err = r.ParseForm()
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // Update recipe with form data
  updatedRecipe := recipe.Recipe{
    Id:        existingRecipe.Id,
    Rev:       existingRecipe.Rev,
    Name:      r.FormValue("r_name"),
    Yields:    r.FormValue("yields"),
    KeepsFor:  r.FormValue("keeps_for"),
    PrepTime:  r.FormValue("prep_time"),
    TotalTime: r.FormValue("total_time"),
  }

  // Process ingredients
  quantities := r.Form["i_quantity[]"]
  names := r.Form["i_name[]"]
  preps := r.Form["i_prep[]"]
  notes := r.Form["i_note[]"]

  for i := 0; i < len(quantities) && i < len(names); i++ {
    ingredient := recipe.Ingredient{
      Quantity: getValue(quantities, i),
      Name:     getValue(names, i),
      Prep:     getValue(preps, i),
      Note:     getValue(notes, i),
    }
    updatedRecipe.Ingredients = append(updatedRecipe.Ingredients, ingredient)
  }

  // Process instructions
  steps := r.Form["step[]"]
  for _, step := range steps {
    if step != "" {
      updatedRecipe.Instructions = append(updatedRecipe.Instructions, step)
    }
  }

  // Process notes
  recipeNotes := r.Form["note[]"]
  for _, note := range recipeNotes {
    if note != "" {
      updatedRecipe.Notes = append(updatedRecipe.Notes, note)
    }
  }

  // Process tags
  tags := r.Form["tag[]"]
  for _, tag := range tags {
    if tag != "" {
      updatedRecipe.Tags = append(updatedRecipe.Tags, tag)
    }
  }

  // Process source
  if r.FormValue("s_name") != "" {
    sourceType := recipe.Type(r.FormValue("s_type"))
    updatedRecipe.Source = recipe.Source{
      Name: r.FormValue("s_name"),
      Type: sourceType,
      Link: r.FormValue("s_link"),
    }
  }

  // Save updated recipe to database
  _, err = db.Save(&updatedRecipe, recipeId, rev)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // Show success page
  templateBytes := appTemplate.MustAsset("web/template/edit_recipe.html")
  t := template.New("edit_recipe.html")
  
  // Add template function for adding numbers
  funcMap := template.FuncMap{
    "add": func(a, b int) int {
      return a + b
    },
  }
  t = t.Funcs(funcMap)
  
  t.Parse(string(templateBytes[:]))
  
  // Create success data with recipe ID for redirect
  successData := map[string]interface{}{
    "Success": true,
    "Id":      recipeId,
  }
  t.Execute(w, successData)
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request, ps httprouter.Params, db *couchdb.Database) {
  recipeId := ps.ByName("id")

  if len(recipeId) == 0 {
    err := errors.New("Recipe ID must be specified")
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // Log the delete attempt
  println("Attempting to delete recipe:", recipeId)

  // Get the recipe first to get the revision
  recipe := recipe.Recipe{}
  rev, err := db.Read(recipeId, &recipe, nil)
  if err != nil {
    println("Error reading recipe:", err.Error())
    http.Error(w, err.Error(), http.StatusNotFound)
    return
  }

  println("Found recipe, revision:", rev)

  // Delete the recipe
  _, err = db.Delete(recipeId, rev)
  if err != nil {
    println("Error deleting recipe:", err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  println("Recipe deleted successfully")

  // Return success response
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Recipe deleted successfully"))
}

func getValue(slice []string, index int) string {
  if index < len(slice) {
    return slice[index]
  }
  return ""
}
