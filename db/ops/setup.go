// Creates the required CouchDB databases and their views.
package ops

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/mryachanin/cookbook/db"
  "github.com/rhinoman/couchdb-go"
  "log"
)

var requiredCouchDbDatabases = [3]string{ "_users", "_replicator", "_global_changes" }

type DesignDocument struct {
  Language string          `json:"language"`
  Views    map[string]View `json:"views"`
}

type View struct {
  Map    string `json:"map"`
  Reduce string `json:"reduce,omitempty"`
}

func SetupDatabase() {
  c := config.LoadConfiguration(configPath)
  conn := db.EstablishConnection(c)
  auth := db.CreateAuth(c)

  // Do nothing if database exists.
  d := conn.SelectDB(db.DatabaseName, auth)
  if err := d.DbExists(); err == nil {
    log.Printf("Database already exists: %s. Exiting", db.DatabaseName)
    return
  }

  // Create required databases.
  createDatabases(conn, auth)

  // Connect to the database.
  d = conn.SelectDB(db.DatabaseName, auth)
  log.Printf("Connected to database \"%s\"", db.DatabaseName)

  // Set up views.
  createViews(d)
}

func createDatabases(conn *couchdb.Connection, auth *couchdb.BasicAuth) {
  // Create the dbs required by CouchDB.
  for _, dbName := range requiredCouchDbDatabases {
    createDatabase(dbName, conn, auth)
  }

  // Create the main app db.
  createDatabase(db.DatabaseName, conn, auth)
}

func createDatabase(dbName string, conn *couchdb.Connection, auth *couchdb.BasicAuth) {
  if err := conn.CreateDB(dbName, auth); err != nil {
    log.Fatalf("Could not create database: %s. Error: %s", dbName, err)
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
