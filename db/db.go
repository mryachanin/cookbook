package db

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/rhinoman/couchdb-go"
  "log"
)

const (
  // Uppercase characters are not allowed.
  DatabaseName = "cookbook"
)

func Connect(c *config.Config) (*couchdb.Database) {
  conn := EstablishConnection(c)
  auth := CreateAuth(c)

  db := conn.SelectDB(DatabaseName, auth)
  if err := db.DbExists(); err != nil {
    log.Fatalf("Failed to connect to database \"%s\". Error: %s", DatabaseName, err)
  }
  log.Printf("Connected to database \"%s\"", DatabaseName)

  return db
}
