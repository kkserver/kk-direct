package main

import (
	"github.com/kkserver/kk-direct/direct/httpd"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	log.SetFlags(log.Llongfile | log.LstdFlags)

	env := "./config/env.ini"

	if len(os.Args) > 1 {
		env = os.Args[1]
	}

	a := httpd.App{}

	err := app.Load(&a, "./app.ini")

	if err != nil {
		log.Panicln(err)
	}

	err = app.Load(&a, env)

	if err != nil {
		log.Panicln(err)
	}

	app.Obtain(&a)

	app.Handle(&a, &app.InitTask{})

	httpd.Openlib()

	go func() {

		http.Handle("/", &a)

		log.Println("httpd " + a.Address)

		srv := &http.Server{
			ReadTimeout:  20 * time.Second,
			WriteTimeout: 30 * time.Second,
			Addr:         a.Address,
		}

		log.Println(srv.ListenAndServe())

	}()

	kk.DispatchMain()

	app.Recycle(&a)

}
