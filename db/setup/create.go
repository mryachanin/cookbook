// Creates the Satisfied Vegan CouchDB database and its views.
package main

import (
  "gitlab.com/mryachanin/satisfied-vegan/config"
  "gitlab.com/mryachanin/satisfied-vegan/db"
  "github.com/rhinoman/couchdb-go"
  "log"
)

func main() {
  c := config.LoadConfiguration()
  conn := db.EstablishConnection(c)
  auth := db.CreateAuth(c)

  // Create the database.
  createDatabase(conn, auth)

  // Connect to the database.
  conn.SelectDB(db.DatabaseName, auth)
  log.Printf("Connected to database \"%s\"", db.DatabaseName)

  // TODO: Set up database views.
}

// Creates the Satisfied Vegan database.
func createDatabase(conn *couchdb.Connection, auth *couchdb.BasicAuth) {
  if err := conn.CreateDB(db.DatabaseName, auth); err != nil {
    log.Fatalf("Could not create database. Error: %s", err)
  }
}
