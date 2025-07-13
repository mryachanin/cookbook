// Creates the required CouchDB databases and their views.
package db

import (
  "github.com/mryachanin/cookbook/config"
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

func SetupDatabase(c *config.Config, conn *couchdb.Connection, auth *couchdb.BasicAuth) {
  // Check if database exists
  d := conn.SelectDB(DatabaseName, auth)
  dbExists := d.DbExists() == nil
  
  if dbExists {
    log.Printf("Database '%s' already exists, checking views...", DatabaseName)
    // Database exists, but we should still check and create views if needed
    createViews(d)
    log.Printf("Database setup verification completed for '%s'", DatabaseName)
    return
  }

  log.Printf("Database '%s' does not exist, creating...", DatabaseName)

  // Create required databases.
  createDatabases(conn, auth)

  // Connect to the database.
  d = conn.SelectDB(DatabaseName, auth)
  log.Printf("Connected to database \"%s\"", DatabaseName)

  // Set up views.
  createViews(d)
  log.Printf("Database setup completed for '%s'", DatabaseName)
}

func createDatabases(conn *couchdb.Connection, auth *couchdb.BasicAuth) {
  // Create the dbs required by CouchDB.
  for _, dbName := range requiredCouchDbDatabases {
    createDatabase(dbName, conn, auth)
  }

  // Create the main app db.
  createDatabase(DatabaseName, conn, auth)
}

func createDatabase(dbName string, conn *couchdb.Connection, auth *couchdb.BasicAuth) {
  log.Printf("Creating database: %s", dbName)

  if err := conn.CreateDB(dbName, auth); err != nil {
    // Check if error is due to database already existing
    if couchErr, ok := err.(*couchdb.Error); ok && couchErr.StatusCode == 412 {
      log.Printf("Database '%s' already exists, skipping creation", dbName)
      return
    }
    log.Fatalf("Could not create database: %s. Error: %s", dbName, err)
  } else {
    log.Printf("Successfully created database: %s", dbName)
  }
}

func createViews(d *couchdb.Database) {
  // Check if design document already exists
  var existingDoc DesignDocument
  rev, err := d.Read("_design/recipe", &existingDoc, nil)
  
  // If document exists, check if we need to update it
  if err == nil {
    log.Printf("Design document 'recipe' already exists, checking if update is needed")
    
    // Check if new view exists
    if _, hasNewView := existingDoc.Views["getRecipesByTag"]; !hasNewView {
      log.Printf("Missing 'getRecipesByTag' view, updating design document")
      // Continue to update the design document
    } else {
      log.Printf("All views exist, skipping creation")
      return
    }
  } else {
    // If error is not "not found", it's a real error
    if couchErr, ok := err.(*couchdb.Error); ok && couchErr.StatusCode != 404 {
      log.Printf("Error checking existing design document: %s", err)
    }
    rev = "" // No existing document
  }

  // Map from view name -> view
  var recipeViews = map[string]View{
    "getRecipes": {
      Map: `
        function(doc) {
          emit(doc.name, doc._id);
        }`,
    },
    "getRecipesByTag": {
      Map: `
        function(doc) {
          if (doc.tags && doc.tags.length > 0) {
            for (var i = 0; i < doc.tags.length; i++) {
              emit(doc.tags[i], {id: doc._id, name: doc.name});
            }
          }
        }`,
    },
  }

  ddoc := DesignDocument{
    Language: "javascript",
    Views:    recipeViews,
  }

  log.Printf("Creating/updating views for design document 'recipe'")

  if _, err := d.SaveDesignDoc("recipe", ddoc, rev); err != nil {
    // Check if error is due to design document already existing (race condition)
    if couchErr, ok := err.(*couchdb.Error); ok && couchErr.StatusCode == 409 {
      log.Printf("Design document 'recipe' was created concurrently, this is normal")
      return
    }
    log.Fatalf("Could not create view. Error: %s", err)
  } else {
    log.Printf("Successfully created views for design document 'recipe'")
  }
}
