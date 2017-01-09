package kk

import (
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/client"
	Value "github.com/kkserver/kk-lib/kk/value"
	"reflect"
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

			opt := reflect.ValueOf(options)

			task := client.RequestTask{}

			task.Name = options.Name()
			task.Timeout = time.Duration(Value.IntValue(Value.Get(opt, "timeout"), 1)) * time.Second

			data := map[string]interface{}{}

			task.Request = data

			v, ok = options["options"]

			if ok {

				mdata, ok := v.(direct.Options)

				if ok {
					for key, value := range mdata {
						vv := direct.ReflectValue(D.App(), ctx, value)
						skey := key.(string)
						if key == "_" {
							Value.EachObject(reflect.ValueOf(vv), func(key reflect.Value, vv reflect.Value) bool {
								if vv.IsValid() && vv.CanInterface() && !vv.IsNil() {
									data[Value.StringValue(key, "")] = vv.Interface()
								}
								return true
							})
						} else {
							data[skey] = vv
						}
					}
				}

			}

			err := app.Handle(a, &task)

			if err != nil {
				return D.Fail(ctx, err)
			}

			ctx.Set(ResultKeys, task.Result)

			return D.Done(ctx, "done")
		}
	}

	return D.Fail(ctx, direct.NewError(direct.ERROR_UNKNOWN, "Not Found kk app"))
}
