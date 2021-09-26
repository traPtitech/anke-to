package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/traPtitech/anke-to/model"
	"github.com/traPtitech/anke-to/tuning"
)

func main() {
	env, ok := os.LookupEnv("ANKE-TO_ENV")
	if !ok {
		env = "production"
	}
	logOn := env == "pprof" || env == "dev"

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			tuning.Inititial()
			return
		case "bench":
			tuning.Bench()
			return
		}
	}

	err := model.EstablishConnection(logOn)
	if err != nil {
		panic(err)
	}

	err = model.Migrate()
	if err != nil {
		panic(err)
	}

	if env == "pprof" {
		runtime.SetBlockProfileRate(1)
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("no PORT")
	}

	SetRouting(port)
}
