package db

import (
  "gitlab.com/mryachanin/satisfied-vegan/config"
  "github.com/rhinoman/couchdb-go"
  "log"
)

const (
  // Uppercase characters are not allowed.
  DatabaseName = "satisfied_vegan"
)

func Connect(c *config.Config) (*couchdb.Database) {
  conn := EstablishConnection(c)
  auth := CreateAuth(c)

  db := conn.SelectDB(DatabaseName, auth)
  if err := db.DbExists(); err != nil {
    log.Fatalf("Failed to connect to database \"%s\"", DatabaseName)
  }
  log.Printf("Connected to database \"%s\"", DatabaseName)

  return db
}
