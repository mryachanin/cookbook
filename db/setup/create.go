// Creates the CouchDB database and its views.
package main

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/mryachanin/cookbook/db"
  "github.com/rhinoman/couchdb-go"
  "log"
)

type DesignDocument struct {
  Language string          `json:"language"`
  Views    map[string]View `json:"views"`
}

type View struct {
  Map    string `json:"map"`
  Reduce string `json:"reduce,omitempty"`
}

func main() {
  c := config.LoadConfiguration()
  conn := db.EstablishConnection(c)
  auth := db.CreateAuth(c)

  // Create the database.
  createDatabase(conn, auth)

  // Connect to the database.
  d := conn.SelectDB(db.DatabaseName, auth)
  log.Printf("Connected to database \"%s\"", db.DatabaseName)

  // Set up views.
  createViews(d)
}

// Create the database.
func createDatabase(conn *couchdb.Connection, auth *couchdb.BasicAuth) {
  if err := conn.CreateDB(db.DatabaseName, auth); err != nil {
    log.Fatalf("Could not create database. Error: %s", err)
  }
}

func createViews(d *couchdb.Database) {
  // Map from view name -> view
  var recipeViews = map[string]View{
    "getRecipes": {
      Map: `
        function(doc) {
          emit(doc.name, doc._id);
        }`,
    },
  }
  ddoc := DesignDocument{
    Language: "javascript",
    Views:    recipeViews,
  }
  if _, err := d.SaveDesignDoc("recipe", ddoc, ""); err != nil {
    log.Fatalf("Could not create view. Error: %s", err)
  }
}
