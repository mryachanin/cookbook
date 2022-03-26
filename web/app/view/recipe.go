package view

import (
  "github.com/mryachanin/cookbook/api/recipe"
  appTemplate "github.com/mryachanin/cookbook/web/template"
  "github.com/rhinoman/couchdb-go"
  "errors"
  "html/template"
  "net/http"
  "net/url"
)

func GetRecipe(w http.ResponseWriter, r *http.Request, db *couchdb.Database) {
  u, err := url.Parse(r.RequestURI)
  if err != nil {
    err := errors.New("Encountered error when parsing URL: " + r.RequestURI)
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  q := u.Query()
  recipeIdQuery := q["id"]

  if recipeIdQuery == nil || len(recipeIdQuery) == 0 {
    err := errors.New("Recipe ID must be specified")
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  recipeId := recipeIdQuery[0]
  recipe := recipe.Recipe{}
  db.Read(recipeId, &recipe, nil)

  templateArr := appTemplate.MustAsset("web/template/recipe.html")
  t := template.New("recipe.html")
  t.Parse(string(templateArr[:]))
  t.Execute(w, recipe)
}
