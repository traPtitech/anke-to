package test

import (
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/router"
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
