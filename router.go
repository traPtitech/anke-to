package main

import (
    "net/http"
    "time"
    "strconv"

    "database/sql"

    "github.com/labstack/echo"

    "github.com/go-sql-driver/mysql"
)

type questionnaires struct {
    ID              int             `json:"questionnaireID" db:"id"`
    Title           string          `json:"title"           db:"title"`
    Description     string          `json:"description"     db:"description"`
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

func timeConvert(time mysql.NullTime) string {
    if time.Valid {
        return time.Time.String()
    } else {
        return "NULL"
    }
}

func getID(c echo.Context) error {
    user := c.Request().Header.Get("X-Showcase-User")

    return c.String(http.StatusOK, "traQID:" + user);
}

// echoに追加するハンドラーは型に注意
// echo.Contextを引数にとってerrorを返り値とする
func getQuestionnaires(c echo.Context) error {
    // query parametar
    sort := c.QueryParam("sort")
    page := c.QueryParam("page")
    //nontargeted := c.QueryParam("nontargeted") == "true"

    if page == "" {
        page = "1"
    }
    num, _ := strconv.Atoi(page)

    var list = map[string] string {
        "":             "",
        "created_at":   "ORDER BY created_at",
        "-created_at":  "ORDER BY created_at DESC",
        "title":        "ORDER BY title",
        "-title":       "ORDER BY title DESC",
        "modified_at":  "ORDER BY modified_at",
        "-modified_at": "ORDER BY modified_at DESC",
    }
    // アンケート一覧の配列
    allquestionnaires := []questionnaires{}

    // これで一気に取れる
    err := db.Select(&allquestionnaires, "SELECT * FROM questionnaires " + list[sort] + " lIMIT 20 OFFSET " + strconv.Itoa(20 * (num-1)))

    // エラー処理
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }
    
    type questionnairesInfo struct {
        ID              int             `json:"questionnaireID"`
        Title           string          `json:"title"`
        Description     string          `json:"description"`
        Res_time_limit  string          `json:"res_time_limit"`
        Deleted_at      string          `json:"deleted_at"`
        Res_shared_to   string          `json:"res_shared_to"`
        Created_at      time.Time       `json:"created_at"`
        Modified_at     time.Time       `json:"modified_at"`
        Is_targeted     bool            `json:"is_targeted"`
    }
    var ret []questionnairesInfo

    for _, v := range allquestionnaires {
        ret = append(ret,
            questionnairesInfo{
                v.ID,
                v.Title,
                v.Description,
                timeConvert(v.Res_time_limit),
                timeConvert(v.Deleted_at),
                v.Res_shared_to,
                v.Created_at,
                v.Modified_at,
                // とりあえず仮でtrueにしている
                true })
    }

    // 構造体の定義で書いたJSONのキーで変換される
    return c.JSON(http.StatusOK, ret)
}

func getQuestionnaire(c echo.Context) error {

    questionnaireID := c.Param("id")
    
    questionnaire := questionnaires{}
    if err := db.Get(&questionnaire, "SELECT * FROM questionnaires WHERE id = ?", questionnaireID);
    err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }
    
    targets := []string{}
    if err := db.Select(&targets, "SELECT user_traqid FROM targets WHERE questionnaire_id = ?", questionnaireID);
    err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    administrators := []string{}
    if err := db.Select(&administrators, "SELECT user_traqid FROM administrators WHERE questionnaire_id = ?", questionnaireID);
    err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    respondents := []string{}
    if err := db.Select(&respondents, "SELECT user_traqid FROM respondents WHERE questionnaire_id = ?", questionnaireID);
    err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "questionnaireID":  questionnaire.ID,
        "title":            questionnaire.Title,
        "description":      questionnaire.Description,
        "res_time_limit":   timeConvert(questionnaire.Res_time_limit),
        "deleted_at":       timeConvert(questionnaire.Deleted_at),
        "created_at":       questionnaire.Created_at,
        "modified_at":      questionnaire.Modified_at,
        "res_shared_to":    questionnaire.Res_shared_to,
        "targets":          targets,
        "administrators":   administrators,
        "respondents":      respondents,
    })
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
        Description     string      `json:"description"`
        Res_time_limit  time.Time   `json:"res_time_limit"`
        Res_shared_to   string      `json:"res_shared_to"`
        Targets         []string    `json:"targets"`
        Administrators  []string    `json:"administrators"`
    }{}

    // JSONを構造体につける
    err := c.Bind(&req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    if req.Title == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "title is null"})
    }
    
    if req.Res_shared_to == "" {
        req.Res_shared_to = "administrators"
    }

    var result sql.Result

    // アンケートの追加    
    if req.Res_time_limit.IsZero() {
        result = db.MustExec(
            "INSERT INTO questionnaires (title, description, res_shared_to) VALUES (?, ?, ?)", 
            req.Title, req.Description, req.Res_shared_to)
    } else {
        result = db.MustExec(
            "INSERT INTO questionnaires (title, description, res_time_limit, res_shared_to) VALUES (?, ?, ?, ?)", 
            req.Title, req.Description, req.Res_time_limit, req.Res_shared_to)
    }

    // エラーチェック
    lastID, err := result.LastInsertId()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    for _, v := range req.Targets {
        _, err := db.Exec(
            "INSERT INTO targets (questionnaire_id, user_traqid) VALUES (?, ?)",
            lastID, v)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }
    }

    for _, v := range req.Administrators {
        _, err := db.Exec(
            "INSERT INTO administrators (questionnaire_id, user_traqid) VALUES (?, ?)",
            lastID, v)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }
    }

    return c.JSON(http.StatusCreated, map[string]interface{}{
        "questionnaireID":  int(lastID),
        "title":            req.Title,
        "description":      req.Description,
        "res_time_limit":   req.Res_time_limit,
        "deleted_at":       "NULL",
        "created_at":       time.Now(),
        "modified_at":      time.Now(),
        "res_shared_to":    req.Res_shared_to,
        "targets":          req.Targets,
        "administrators":   req.Administrators,
    })
}

func editQuestionnaire(c echo.Context) error {
    questionnaireID := c.Param("id")
    
    req := struct {
        Title           string      `json:"title"`
        Description     string      `json:"description"`
        Res_time_limit  time.Time   `json:"res_time_limit"`
        Res_shared_to   string      `json:"res_shared_to"`
        Targets         []string    `json:"targets"`
        Administrators  []string    `json:"administrators"`
    }{}

    err := c.Bind(&req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, err)
    }

    if req.Title == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "title is null"})
    }
    
    if req.Res_shared_to == "" {
        req.Res_shared_to = "administrators"
    }

    // アップデートする
    if req.Res_time_limit.IsZero() {
        if _, err := db.Exec(
            "UPDATE questionnaires SET title = ?, description = ?, res_shared_to = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?", 
            req.Title, req.Description, req.Res_shared_to, questionnaireID);
        err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }
    } else {
        if _, err = db.Exec(
            "UPDATE questionnaires SET title = ?, description = ?, res_time_limit = ?, res_shared_to = ?, modified_at = CURRENT_TIMESTAMP WHERE id = ?", 
            req.Title, req.Description, req.Res_time_limit, req.Res_shared_to, questionnaireID);
        err != nil {
            return c.JSON(http.StatusInternalServerError, err)
        }
    }

    // TargetsとAdministratorsの変更はまだ

    return c.NoContent(http.StatusOK)
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