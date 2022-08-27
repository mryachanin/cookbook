package app

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/mryachanin/cookbook/web/app/view"
  "github.com/julienschmidt/httprouter"
  "github.com/rhinoman/couchdb-go"
  "log"
  "net/http"
  "strconv"
)

func HandleRequests(c *config.Config, db *couchdb.Database) {
  router := httprouter.New()

  router.GET("/", wrap(view.GetIndex, db))
  router.GET("/recipe", wrap(view.CreateRecipe, db))
  router.POST("/recipe", wrap(view.PostRecipe, db))
  router.GET("/recipe/:id", wrap(view.GetRecipe, db))

  url := c.Host + ":" + strconv.Itoa(c.Port)
  log.Fatal(http.ListenAndServe(url, router))
}

func wrap(handler func(http.ResponseWriter, *http.Request, httprouter.Params, *couchdb.Database),
          db *couchdb.Database) (func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)) {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    handler(w, r, ps, db)
  }
}
