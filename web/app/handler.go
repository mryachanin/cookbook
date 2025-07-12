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

  // Static CSS file serving
  router.GET("/styles.css", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    w.Header().Set("Content-Type", "text/css")
    http.ServeFile(w, r, "web/template/styles.css")
  })

  router.GET("/", wrap(view.GetIndex, db))
  router.GET("/create", wrap(view.CreateRecipe, db))
  router.POST("/create", wrap(view.PostRecipe, db))
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
