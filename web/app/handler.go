package app

import (
  "gitlab.com/mryachanin/satisfied-vegan/config"
  "gitlab.com/mryachanin/satisfied-vegan/web/app/view"
  "github.com/rhinoman/couchdb-go"
  "log"
  "net/http"
  "strconv"
)

func HandleRequests(c *config.Config, db *couchdb.Database) {
  http.HandleFunc("/", wrap(view.GetIndex, db))
  http.HandleFunc("/recipe", wrap(view.GetRecipe, db))
  url := c.Host + ":" + strconv.Itoa(c.Port)
  log.Fatal(http.ListenAndServe(url, nil))
}

func wrap(handler func(http.ResponseWriter, *http.Request, *couchdb.Database),
          db *couchdb.Database) (func(w http.ResponseWriter, r *http.Request)) {
  return func(w http.ResponseWriter, r *http.Request) {
    handler(w, r, db)
  }
}
