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

func main(){

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
    e.GET("/questionnaire", getQuestionnaire)
    e.POST("/questionnaire", postQuestionnaire)
    e.PATCH("/questionnaire/:id", editQuestionnaire)
    e.DELETE("/questionnaire/:id", deleteQuestionnaire)
    e.GET("/questionnaire/:id", getQuestions)

    e.GET("/users/me", getID)

    // Start server
    e.Logger.Fatal(e.Start(":1323"))
}