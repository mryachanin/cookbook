package main

import (
  "gitlab.com/mryachanin/satisfied-vegan/config"
  "gitlab.com/mryachanin/satisfied-vegan/db"
  "gitlab.com/mryachanin/satisfied-vegan/web/app"
)

func main() {
  c := config.LoadConfiguration()
  db := db.Connect(c)
  app.HandleRequests(c, db)
}
