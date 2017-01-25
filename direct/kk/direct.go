package kk

import (
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/client"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"log"
	"strings"
	"time"
)

var AppKeys = []string{"app"}
var ResultKeys = []string{"result"}

type Direct struct {
	direct.Direct
}

func (D *Direct) Exec(ctx direct.IContext) error {

	options := D.Options()

	v := ctx.Get(AppKeys)

	if v != nil {

		a, ok := v.(app.IApp)

		if ok {

			task := client.RequestTask{}

			task.Name = options.Name()

			route := ctx.Get(RouteKeys)

			if route != nil {
				dynamic.Each(route, func(key interface{}, value interface{}) bool {
					skey := dynamic.StringValue(key, "")
					log.Println("skey", skey)
					if strings.HasPrefix(task.Name, skey) {
						task.Name = dynamic.StringValue(value, "") + task.Name[len(skey):]
						return false
					}
					return true
				})
			}

			task.Timeout = time.Duration(dynamic.IntValue(dynamic.Get(options, "timeout"), 1)) * time.Second

			v, ok = options["options"]

			if ok {
				task.Request = direct.ReflectValue(D.App(), ctx, v)
			} else {
				task.Request = map[interface{}]interface{}{}
			}

			log.Println("kk.direct", task.Name, task.Request)

			err := app.Handle(a, &task)

			log.Println("kk.direct", task.Name, task.Result, err)

			if err != nil {
				return D.Fail(ctx, err)
			}

			errno := int(dynamic.IntValue(dynamic.Get(task.Result, "errno"), 0))
			errmsg := dynamic.StringValue(dynamic.Get(task.Result, "errmsg"), "")

			if errno != 0 {
				return D.Fail(ctx, direct.NewError(errno, errmsg))
			}

			ctx.Set(ResultKeys, task.Result)

			return D.Done(ctx, "done")
		}
	}

	return D.Fail(ctx, direct.NewError(direct.ERROR_UNKNOWN, "Not Found kk app"))
}
