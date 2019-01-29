package main

import (
  "gitlab.com/mryachanin/satisfied-vegan/web/app"
)

func main() {
  app.HandleRequests()
}

func init() {
	// init db
}
