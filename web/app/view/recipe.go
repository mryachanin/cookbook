package view

import (
  "github.com/mryachanin/cookbook/api/recipe"
  appTemplate "github.com/mryachanin/cookbook/web/template"
  "github.com/julienschmidt/httprouter"
  "github.com/rhinoman/couchdb-go"
  "errors"
  "fmt"
  "html/template"
  "io/ioutil"
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
  bodyBytes, _ := ioutil.ReadAll(r.Body)
  bodyString := string(bodyBytes)
  fmt.Fprintf(w, bodyString)
}
