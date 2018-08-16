package main

import (
    "net/http"
    "time"

    "database/sql"

    "github.com/labstack/echo"

    "github.com/go-sql-driver/mysql"
)

type questionnaires struct {
    ID              int             `json:"questionnaireID" db:"id"`
    Title           string          `json:"title"           db:"title"`
    Res_time_limit  mysql.NullTime  `json:"res_time_limit"  db:"res_time_limit"`
    Deleted_at      mysql.NullTime  `json:"deleted_at"      db:"deleted_at"`
    Res_shared_to   string          `json:"res_shared_to"   db:"res_shared_to"`
    Created_at      time.Time       `json:"created_at"      db:"created_at"`
    Modified_at     time.Time       `json:"modified_at"     db:"modified_at"`
}

type questions struct {
    ID                  int             `json:"id"                  db:"id"`
    Questionnaire_id    int             `json:"questionnaireID"     db:"questionnaire_id"`
    Page_num            int             `json:"page_num"            db:"page_num"`
    Question_num        int             `json:"question_num"        db:"question_num"`
    Type                string          `json:"type"                db:"type"`
    Body                string          `json:"body"                db:"body"`
    Is_requrired        bool            `json:"is_required"         db:"is_required"`
    Deleted_at          mysql.NullTime  `json:"deleted_at"          db:"deleted_at"`
    Created_at          time.Time       `json:"created_at"          db:"created_at"`
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
    req := struct {
        Title           string      `json:"title"`
        Res_time_limit  time.Time   `json:"res_time_limit"`
        Res_shared_to   string      `json:"res_shared_to"`
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

    var result sql.Result
    var res_time_limit mysql.NullTime

    // アンケートの追加    
    if req.Res_time_limit.IsZero() {
        result = db.MustExec(
            "INSERT INTO questionnaires (title, res_shared_to) VALUES (?, ?)", 
            req.Title, req.Res_shared_to)
        res_time_limit.Valid = false
    } else {
        result = db.MustExec(
            "INSERT INTO questionnaires (title, res_time_limit, res_shared_to) VALUES (?, ?, ?)", 
            req.Title, req.Res_time_limit, req.Res_shared_to)
        res_time_limit = mysql.NullTime{req.Res_time_limit, true}
    }

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

    var deleted_at mysql.NullTime
    deleted_at.Valid = false

    t := questionnaires{
        ID:             int(lastID),
        Title:          req.Title,
        Res_time_limit: res_time_limit,
        Deleted_at:     deleted_at,
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
    result := db.MustExec("UPDATE questionnaires SET title = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?", req.Title, questionnaireID)
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

    result := db.MustExec("UPDATE questionnaires SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", questionnaireID)
    _, err := result.LastInsertId()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    return c.NoContent(http.StatusOK)
}