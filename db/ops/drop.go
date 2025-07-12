// This contains logic to drop the CouchDB database
// WARNING: This will permanently delete all data!
package ops

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/mryachanin/cookbook/db"
  "log"
)

func DropDatabase() {
  c := config.LoadConfiguration(configPath)
  conn := db.EstablishConnection(c)
  auth := db.CreateAuth(c)

  // Check if database exists first
  database := conn.SelectDB(db.DatabaseName, auth)
  if err := database.DbExists(); err != nil {
    log.Printf("Database '%s' does not exist or is not accessible", db.DatabaseName)
    return
  }

  log.Printf("WARNING: About to drop database '%s' - this will permanently delete all data!", db.DatabaseName)
  
  // Drop the database
  if err := conn.DeleteDB(db.DatabaseName, auth); err != nil {
    log.Fatalf("Failed to drop database '%s': %v", db.DatabaseName, err)
  }

  log.Printf("Successfully dropped database '%s'", db.DatabaseName)
}