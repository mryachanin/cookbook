// Creates the required CouchDB databases and their views.
package ops

import (
	"github.com/mryachanin/cookbook/config"
	"github.com/mryachanin/cookbook/db"
  )

func SetupDatabase() {
	c := config.LoadConfiguration(configPath)
	conn := db.EstablishConnection(c)
	auth := db.CreateAuth(c)

	db.SetupDatabase(c, conn, auth)
}  