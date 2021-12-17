package model

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/prometheus"
)

var (
	db        *gorm.DB
	allTables = []interface{}{
		Questionnaires{},
		Questions{},
		Respondents{},
		Responses{},
		Administrators{},
		Options{},
		ScaleLabels{},
		Targets{},
		Validations{},
	}
)

// EstablishConnection DBと接続
func EstablishConnection(isProduction bool) error {
	user, ok := os.LookupEnv("MARIADB_USERNAME")
	if !ok {
		user = "root"
	}

	pass, ok := os.LookupEnv("MARIADB_PASSWORD")
	if !ok {
		pass = "password"
	}

	host, ok := os.LookupEnv("MARIADB_HOSTNAME")
	if !ok {
		host = "localhost"
	}
	port, ok := os.LookupEnv("MARIADB_PORT")
	if !ok {
		port = "3306"
	}

	dbname, ok := os.LookupEnv("MARIADB_DATABASE")
	if !ok {
		dbname = "anke-to"
	}

	var logLevel logger.LogLevel
	if isProduction {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname) + "?parseTime=true&loc=Asia%2FTokyo&charset=utf8mb4"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	db = db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci")

	db.Use(prometheus.New(prometheus.Config{
		DBName:          "anke-to",
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&MetricsCollector{},
		},
	}))

	return err
}

// Migrate DBのMigrationを行う
func Migrate() error {
	err := db.AutoMigrate(allTables...)
	if err != nil {
		return fmt.Errorf("failed in table's migration: %w", err)
	}

	return nil
}
