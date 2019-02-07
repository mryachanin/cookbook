package db

import (
  "gitlab.com/mryachanin/satisfied-vegan/config"
  "github.com/rhinoman/couchdb-go"
  "log"
  "time"
)

const (
  connectTimeout = time.Duration(500 * time.Millisecond)
)

func EstablishConnection(c *config.Config) (*couchdb.Connection) {
  host := c.Database.Host
  port := c.Database.Port
  conn, err := couchdb.NewConnection(host, port, connectTimeout)
  if err != nil {
    log.Fatalf("Could not create a connection to the supplied database at host: %s, port: %d",
      host, port)
  }
  return conn
}
