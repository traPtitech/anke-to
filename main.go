package main

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

var (
    db *sqlx.DB
)

type questionnaires struct {
    ID              int         `json:"id"              db:"id"`
    Title           string      `json:"title"           db:"title"`
    Is_accepting    bool        `json:"is_accepting"    db:"is_accepting"`
    Is_deleted      bool        `json:"is_deleted"      db:"is_deleted"`
    Res_shared_to   string      `json:"res_shared_to"   db:"res_shared_to"`
    Created_at      time.Time   `json:"created_at"      db:"created_at"`
    Modified_at     time.Time   `json:"modified_at"     db:"modified_at"`
}

type questions struct {
    ID                  int         `json:"id"                  db:"id"`
    Questionnaire_id    int         `json:"questionnaire_id"    db:"questionnaire_id"`
    Page_num            int         `json:"page_num"            db:"page_num"`
    Question_num        int         `json:"question_num"        db:"question_num"`
    Type                string      `json:"type"                db:"type"`
    Body                string      `json:"body"                db:"body"`
    Is_requrired        bool        `json:"is_required"         db:"is_required"`
    Is_deleted          bool        `json:"is_deleted"          db:"is_deleted"`
    Created_at          time.Time   `json:"created_at"          db:"created_at"`
}

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

func getID(c echo.Context) error {
    user := c.Request().Header.Get("X-Showcase-User")

    return c.String(http.StatusOK, "traQID:" + user);
}

// echoに追加するハンドラーは型に注意
// echo.Contextを引数にとってerrorを返り値とする
func getQuestionnaire(c echo.Context) error {
    // アンケート一覧の配列
    allquestionnaires := []questionnaires{}

    // これで一気に取れる
    err := db.Select(&allquestionnaires, "SELECT * FROM questionnaires")

    // エラー処理
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    // 構造体の定義で書いたJSONのキーで変換される
    return c.JSON(http.StatusOK, allquestionnaires)
}

func getQuestions(c echo.Context) error {
    questionnaireID := c.Param("id")
    // 質問一覧の配列
    allquestions := []questions{}

    // アンケートidの一致する質問を取る
    err := db.Select(&allquestions, "SELECT * FROM questions WHERE questionnaire_id = ?", questionnaireID)

    // エラー処理
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    // 構造体の定義で書いたJSONのキーで変換される
    return c.JSON(http.StatusOK, allquestions)
}

func postQuestionnaire(c echo.Context) error {

    // リクエストで投げられるJSONのスキーマ
    // その場で構造体を定義もできる
    req := struct {
        Title           string  `json:"title"`
        Is_accepting    bool    `json:"is_accepting"`
        Res_shared_to   string  `json:"res_shared_to"`
        Questions []struct {
            Page_num        int     `json:"page_num"`
            Type            string  `json:"type"`
            Body            string  `json:"body"`
            Is_required     bool    `json:"is_required"`
        }
    }{}

    // JSONを構造体につける
    err := c.Bind(&req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    if req.Title == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "name is null"})
    }
    
    if req.Res_shared_to == "" {
        req.Res_shared_to = "administrators"
    }

    // アンケートの追加
    result := db.MustExec(
        "INSERT INTO questionnaires (title, is_accepting, res_shared_to) VALUES (?, ?, ?)", 
        req.Title, req.Is_accepting, req.Res_shared_to)

    // エラーチェック
    lastID, err := result.LastInsertId()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    for i, v := range req.Questions {
        res := db.MustExec(
            "INSERT INTO questions (questionnaire_id, page_num, question_num, type, body, is_required) VALUES (?, ?, ?, ?, ?, ?)",
            lastID, v.Page_num, i+1, v.Type, v.Body, v.Is_required)

        _, err2 := res.LastInsertId()
        if err2 != nil {
            return c.JSON(http.StatusInternalServerError, err2)
        }
    }

    t := questionnaires{
        ID:             int(lastID),
        Title:          req.Title,
        Is_accepting:   req.Is_accepting,
        Is_deleted:     false,
        Res_shared_to:  req.Res_shared_to,
        Created_at:     time.Now(),
        Modified_at:    time.Now(),
    }
    return c.JSON(http.StatusCreated, t)
}

func editQuestionnaire(c echo.Context) error {
    questionnaireID := c.Param("id")
    req := struct {
        Title string `json:"title"`
    }{}
    err := c.Bind(&req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    if req.Title == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "title is null"})
    }

    // アップデートする
    result := db.MustExec("UPDATE questionnaires SET title = ?, modified_at = ? WHERE id = ?", req.Title, time.Now(), questionnaireID)
    _, err = result.LastInsertId()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    t := questionnaires{}
    db.Get(&t, "SELECT * FROM questionnaires WHERE id = ?", questionnaireID)
    return c.JSON(http.StatusOK, t)
}

func deleteQuestionnaire(c echo.Context) error {
    questionnaireID := c.Param("id")

    result := db.MustExec("UPDATE questionnaires SET is_deleted = True WHERE id = ?", questionnaireID)
    _, err := result.LastInsertId()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    return c.NoContent(http.StatusOK)
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
    
    // Routes
    e.GET("/", getID)
    e.GET("/questionnaire", getQuestionnaire)
    e.POST("/questionnaire", postQuestionnaire)
    e.PATCH("/questionnaire/:id", editQuestionnaire)
    e.DELETE("/questionnaire/:id", deleteQuestionnaire)
    e.GET("/questionnaire/:id", getQuestions)

    // Start server
    e.Logger.Fatal(e.Start(":1323"))
}