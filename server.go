package main

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/mryachanin/cookbook/db"
  "github.com/mryachanin/cookbook/web/app"
)

const (
  configPath = "./config.json"
)

func main() {
  c := config.LoadConfiguration(configPath)
  db := db.Connect(c)
  app.HandleRequests(c, db)
}
