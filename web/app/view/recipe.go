package view

import (
  "gitlab.com/mryachanin/satisfied-vegan/api/recipe"
  appTemplate "gitlab.com/mryachanin/satisfied-vegan/web/template"
  "github.com/rhinoman/couchdb-go"
  "html/template"
  "log"
  "net/http"
  "net/url"
)

func GetRecipe(w http.ResponseWriter, r *http.Request, db *couchdb.Database) {
  u, err := url.Parse(r.RequestURI)
  if err != nil {
    log.Fatal(err)
  }

  q := u.Query()
  recipeId := q["id"][0]
  recipe := recipe.Recipe{}
  db.Read(recipeId, &recipe, nil)

  templateArr := appTemplate.MustAsset("web/template/recipe.html")
  t := template.New("recipe.html")
  t.Parse(string(templateArr[:]))
  t.Execute(w, recipe)
}
