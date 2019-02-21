package model

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

var (
	DB *sqlx.DB
)

func EstablishConnection() (*sqlx.DB, error) {
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

	return sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true&loc=Japan&charset=utf8mb4", user, pass, host, dbname))
}

func SetTimeZone() error {
	_, err := DB.Exec("set time_zone = '+09:00'")
	return err
}
