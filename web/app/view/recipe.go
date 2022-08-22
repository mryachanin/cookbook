package view

import (
  "github.com/mryachanin/cookbook/api/recipe"
  appTemplate "github.com/mryachanin/cookbook/web/template"
  "github.com/julienschmidt/httprouter"
  "github.com/rhinoman/couchdb-go"
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

  templateArr := appTemplate.MustAsset("web/template/recipe.html")
  t := template.New("recipe.html")
  t.Parse(string(templateArr[:]))
  t.Execute(w, recipe)
}
