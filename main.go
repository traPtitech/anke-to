package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func establishConnection() (*sqlx.DB, error) {
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

func main() {

	_db, err := establishConnection()
	if err != nil {
		panic(err)
	}
	db = _db

	e := echo.New()
	e.Use(middleware.CORS())

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Static Files
	e.Static("/", "client/dist")

	// Routes
	e.GET("/questionnaires", getQuestionnaires)
	e.POST("/questionnaires", postQuestionnaire)
	e.GET("/questionnaires/:id", getQuestionnaire)
	e.PATCH("/questionnaires/:id", editQuestionnaire)
	e.DELETE("/questionnaires/:id", deleteQuestionnaire)
	e.GET("/questionnaires/:id/questions", getQuestions)

	e.POST("/questions", postQuestion)
	e.PATCH("/questions/:id", editQuestion)
	e.DELETE("/questions/:id", deleteQuestion)

	e.GET("/users/me", getID)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
