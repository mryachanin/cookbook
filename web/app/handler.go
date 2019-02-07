package app

import (
  "gitlab.com/mryachanin/satisfied-vegan/api/recipe"
  "gitlab.com/mryachanin/satisfied-vegan/config"
  appTemplate "gitlab.com/mryachanin/satisfied-vegan/web/template"
  "github.com/rhinoman/couchdb-go"
  "html/template"
  "log"
  "net/http"
  "strconv"
)

func HandleRequests(c *config.Config, db *couchdb.Database) {
  http.HandleFunc("/", wrap(handleRequest, db))
  url := c.Host + ":" + strconv.Itoa(c.Port)
  log.Fatal(http.ListenAndServe(url, nil))
}

func wrap(handler func(http.ResponseWriter, *http.Request, *couchdb.Database),
          db *couchdb.Database) (func(w http.ResponseWriter, r *http.Request)) {
  return func(w http.ResponseWriter, r *http.Request) {
    handler(w, r, db)
  }
}

func handleRequest(w http.ResponseWriter, r *http.Request, db *couchdb.Database) {
  recipes, err := getRecipes(db)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }
  templateArr := appTemplate.MustAsset("web/template/index.html")
  t := template.New("index.html")
  t.Parse(string(templateArr[:]))
  t.Execute(w, recipes)
}

func getRecipes(db *couchdb.Database) ([]recipe.Recipe, error) {
  r1 := recipe.Recipe{}
  if _, err := db.Read("Mushroom Lentil Loaf", &r1, nil); err != nil {
    log.Fatalf("Could not load recipes. Error: %s", err)
  }
  r2:= recipe.Recipe{}
  if _, err := db.Read("Vegan \"No Tuna\" Salad Sandwich", &r2, nil); err != nil {
    log.Fatalf("Could not load recipes. Error: %s", err)
  }
  return []recipe.Recipe{r1, r2}, nil
}
