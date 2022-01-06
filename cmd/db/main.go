// Creates the CouchDB database and its views.
package main

import (
  "flag"
  "github.com/mryachanin/cookbook/db/ops"
  "log"
)

func main() {
  var op string
  var usage bool

	flag.StringVar(&op, "o", "", "Operation to perform. Valid choices: setup")
	flag.BoolVar(&usage, "h", false, "Display usage info.")

	flag.Parse()

	if op != "" {
		performOp(op)
	} else if usage {
    flag.Usage()
  } else {
    log.Fatalf("No valid command detected. Print usage using -h. Exiting.")
  }
}

func performOp(opName string) {
	if opName == "setup" {
    ops.SetupDatabase()
  } else {
    log.Fatalf("No valid operation detected with name: %s", opName)
  }
}
