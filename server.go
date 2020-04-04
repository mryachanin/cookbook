package main

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/mryachanin/cookbook/db"
  "github.com/mryachanin/cookbook/web/app"
)

func main() {
  c := config.LoadConfiguration()
  db := db.Connect(c)
  app.HandleRequests(c, db)
}
