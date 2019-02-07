package db

import (
  "gitlab.com/mryachanin/satisfied-vegan/config"
  "github.com/rhinoman/couchdb-go"
  "log"
  "time"
)

const (
  name = "satisfied_vegan"  // Uppercase characters are not allowed.
  timeout = time.Duration(500 * time.Millisecond)
)

func Connect(c *config.Config) (*couchdb.Database) {
  conn := establishConnection(c)

  user := c.Database.User
  password := c.Database.Password
  auth := &couchdb.BasicAuth{Username: user, Password: password}

  // Create database if it does not exist.
  if (!isDatabaseCreated(conn)) {
    log.Printf("Database \"%s\" was not found. Creating it now", name)
    createDatabase(conn, auth)
  } else {
    log.Printf("Database \"%s\" was found", name)
  }

  db := conn.SelectDB(name, auth)
  log.Printf("Connected to database \"%s\"", name)

  return db
}

func establishConnection(c *config.Config) (*couchdb.Connection) {
  host := c.Database.Host
  port := c.Database.Port
  conn, err := couchdb.NewConnection(host, port, timeout)
  if err != nil {
    log.Fatalf("Could not create a connection to the supplied database at host: %s, port: %d",
      host, port)
  }
  return conn
}

func isDatabaseCreated(conn *couchdb.Connection) (bool) {
  dbList, err := conn.GetDBList()

  if err != nil {
    log.Fatalf("Could not list databases. Error: %s", err)
  }

  for _, db := range dbList {
    if db == name {
      return true
    }
  }
  return false
}

func createDatabase(conn *couchdb.Connection, auth *couchdb.BasicAuth) {
  if err := conn.CreateDB(name, auth); err != nil {
    log.Fatalf("Could not create database. Error: %s", err)
  }
}
