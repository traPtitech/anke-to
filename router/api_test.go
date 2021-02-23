package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
)

type users string
type httpMethods string
type contentTypes string

const (
	rootPath               = "/api"
	userHeader             = "X-Showcase-User"
	userUnAuthorized       = "-"
	userOne          users = "mazrean"
	userTwo          users = "ryoha"
	//userThree        users        = "YumizSui"
	methodGet    httpMethods  = http.MethodGet
	methodPost   httpMethods  = http.MethodPost
	methodPatch  httpMethods  = http.MethodPatch
	methodDelete httpMethods  = http.MethodDelete
	typeNone     contentTypes = ""
	typeJSON     contentTypes = echo.MIMEApplicationJSON
)

var (
	errMock = errors.New("Mock Error")
)

func makePath(path string) string {
	return rootPath + path
}

func createRecorder(e *echo.Echo, user users, method httpMethods, path string, contentType contentTypes, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(string(method), path, strings.NewReader(body))
	if contentType != typeNone {
		req.Header.Set(echo.HeaderContentType, string(contentType))
	}
	req.Header.Set(userHeader, string(user))

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec
}
