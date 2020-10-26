package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/router"
)

type users string
type httpMethods string
type contentTypes string

const (
	rootPath = "/api"
	userHeader = "X-Showcase-User"
	userUnAuthorized = "-"
	userOne users = "mazrean"
	userTwo users = "ryoha"
	userThree users = "YumizSui"
	methodGet httpMethods = http.MethodGet
	methodPost httpMethods = http.MethodPost
	methodPatch httpMethods = http.MethodPatch
	methodDelet httpMethods = http.MethodDelete
	typeNone contentTypes = ""
	typeJSON contentTypes = echo.MIMEApplicationJSON
)

var e *echo.Echo

//TestMain テストのmain
func TestMain(m *testing.M) {
	db, err := model.EstablishConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = model.Migrate()
	if err != nil {
		panic(err)
	}

	e = echo.New()
	router.SetRouting(e)

	code := m.Run()

	os.Exit(code)
}

func createRecorder(user users, method httpMethods, path string, contentType contentTypes, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(string(method), path, strings.NewReader(body))
	if contentType != typeNone {
		req.Header.Set(echo.HeaderContentType, string(contentType))
	}
	req.Header.Set(userHeader, string(user))

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}

func makePath(path string) string {
	return rootPath + path
}
