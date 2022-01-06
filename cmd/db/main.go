// Creates the CouchDB database and its views.
package main

import (
  "flag"
  "github.com/mryachanin/cookbook/db/ops"
  "log"
)

func main() {
  var op string
  var importPath string
  var usage bool

	flag.StringVar(&op, "o", "", "Operation to perform. Valid choices: import, setup")
  flag.StringVar(&importPath, "i", "",
    "A path to a directory containing recipes in YAML to load into CouchDB.\n" +
    "This is required when passing -o import")
	flag.BoolVar(&usage, "h", false, "Display usage info.")

	flag.Parse()

	if usage {
    flag.Usage()
  } else if op == "setup" {
    ops.SetupDatabase()
  } else if op == "import" {
    if importPath == "" {
      log.Fatal("Import path must be specified using the '-i' flag. Exiting.")
    }
    ops.ImportRecipes(importPath)
  } else {
    log.Fatalf("No valid command detected. Print usage using -h. Exiting.")
  }
}
