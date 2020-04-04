package db

import (
  "github.com/mryachanin/cookbook/config"
  "github.com/rhinoman/couchdb-go"
)

func CreateAuth(c *config.Config) (*couchdb.BasicAuth) {
  user := c.Database.User
  password := c.Database.Password
  return &couchdb.BasicAuth{Username: user, Password: password}
}
