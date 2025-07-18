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
  router.POST("/recipe/:id", wrap(view.UpdateRecipe, db))
  router.DELETE("/recipe/:id", wrap(view.DeleteRecipe, db))
  router.POST("/recipe/:id/delete", wrap(view.DeleteRecipe, db))
  router.GET("/edit/:id", wrap(view.EditRecipe, db))
  router.POST("/edit/:id", wrap(view.UpdateRecipe, db))
  router.GET("/tag/:tag", wrap(view.GetRecipesByTag, db))
  
  // Shopping list routes
  router.GET("/shopping", wrap(view.GetShoppingList, db))
  router.POST("/recipe/:id/add-to-shopping", wrap(view.AddToShoppingList, db))
  router.POST("/shopping/remove", wrap(view.RemoveFromShoppingList, db))
  router.POST("/shopping/toggle", wrap(view.ToggleShoppingItem, db))

  url := c.Host + ":" + strconv.Itoa(c.Port)
  log.Fatal(http.ListenAndServe(url, router))
}

func wrap(handler func(http.ResponseWriter, *http.Request, httprouter.Params, *couchdb.Database),
          db *couchdb.Database) (func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)) {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    handler(w, r, ps, db)
  }
}
