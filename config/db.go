package config

import (
	"os"

	"github.com/gocroot/helper/atdb"
)

var MongoString string = os.Getenv("MONGOSTRINGGIS")

var mongoinfo = atdb.DBInfo{
	DBString: MongoString,
	DBName:   "GIS",
}

var Mongoconn, ErrorMongoconn = atdb.MongoConnect(mongoinfo)
