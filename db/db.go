package db

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/rhinoman/couchdb-go"
  "log"
  "math"
  "time"
)

const (
  // Uppercase characters are not allowed.
  DatabaseName = "cookbook"
  retries = 3
  setupDbIfNotExists = true
)

func Connect(c *config.Config) (*couchdb.Database) {
  conn := EstablishConnection(c)
  auth := CreateAuth(c)

  db := connectWithRetry(c, conn, auth)

  log.Printf("Connected to database \"%s\"", DatabaseName)

  return db
}

func connectWithRetry(c *config.Config, conn *couchdb.Connection, auth *couchdb.BasicAuth) (*couchdb.Database) {
  for i:=0; i<retries; i++ {
    if i > 0 {
      backoff := i * 5 * int(math.Pow(10, float64(i-1)))

      log.Printf("Retrying connection in %d ms")

      time.Sleep(time.Duration(backoff) * time.Millisecond)
    }

    db := conn.SelectDB(DatabaseName, auth)

    err := db.DbExists()
    if err == nil {
      return db
    }

    log.Printf("Failed to connect to database \"%s\". Error: %s", DatabaseName, err)

    if setupDbIfNotExists {
      SetupDatabase(c, conn, auth)
    }
  }

  log.Fatal("Reached max number of failures while connecting to database. Exiting.")
  return nil
}
