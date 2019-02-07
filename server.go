package main

import (
  "gitlab.com/mryachanin/satisfied-vegan/config"
  "gitlab.com/mryachanin/satisfied-vegan/db"
  "gitlab.com/mryachanin/satisfied-vegan/web/app"
)

const (
  configPath = "./config.json"
)

func main() {
  c := config.LoadConfiguration(configPath)
  db := db.Connect(c)
  app.HandleRequests(c, db)
}
