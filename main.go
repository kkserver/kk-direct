package main

import (
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-direct/direct/lua"
	"github.com/kkserver/kk-direct/direct/yaml"
	"log"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	direct.Openlib()
	yaml.Openlib()
	lua.Openlib()

	app, err := yaml.Load("./app.yaml")

	if err != nil {
		log.Panicln(err)
	}

	ctx := direct.NewContext()

	ctx.Begin()

	lua.ContextOpenlib(ctx)

	err = app.Exec(ctx)

	if err != nil {
		log.Panicln(err)
	}

	if err != nil {
		log.Panicln(err)
	}

	log.Println(ctx.Get(direct.OutputKeys))

	lua.ContextCloselib(ctx)

	ctx.End()

}
