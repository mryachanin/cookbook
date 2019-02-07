package main

import (
  "gitlab.com/mryachanin/satisfied-vegan/config"
  "gitlab.com/mryachanin/satisfied-vegan/web/app"
)

const (
  configPath = "./config.json"
)

func main() {
  c := config.LoadConfiguration(configPath)
  app.HandleRequests(c)
}
