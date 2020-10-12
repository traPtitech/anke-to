package model

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var gormDB *gorm.DB

func EstablishConnection() (*gorm.DB, error) {
	user := os.Getenv("MARIADB_USERNAME")
	if user == "" {
		user = "root"
	}

	pass := os.Getenv("MARIADB_PASSWORD")
	if pass == "" {
		pass = "password"
	}

	host := os.Getenv("MARIADB_HOSTNAME")
	if host == "" {
		host = "localhost"
	}

	dbname := os.Getenv("MARIADB_DATABASE")
	if dbname == "" {
		dbname = "anke-to"
	}

	_db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4")
	gormDB = _db
	gormDB = gormDB.BlockGlobalUpdate(true)
	db = sqlx.NewDb(_db.DB(), "mysql")

	return gormDB, err
}
