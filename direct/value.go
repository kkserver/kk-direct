package direct

import (
	"github.com/kkserver/kk-lib/kk/dynamic"
	"reflect"
	"strings"
)

type ValueNil struct {
}

var Nil = &ValueNil{}

var ResultKeys = []string{"result"}

func ReflectValue(app IApp, ctx IContext, value interface{}) interface{} {

	if value != nil {

		{
			v, ok := value.(string)

			if ok {
				if strings.HasPrefix(v, "=") {
					return ctx.Get(strings.Split(v[1:], "."))
				} else if strings.HasPrefix(v, "^") {
					idx := strings.Index(v, " ")
					options := Options{}
					if idx >= 0 {
						options["name"] = v[:idx]
						options["content"] = v[idx+1:]
					} else {
						options["name"] = v
					}
					v, err := app.Open(options)
					if err != nil {
						return err
					}
					ctx.Begin()
					ctx.Set(ResultKeys, Nil)
					err = v.Exec(ctx)
					if err != nil {
						value = err
					} else {
						value = ctx.Get(ResultKeys)
					}
					ctx.End()
					return value
				}
			}
		}

		{
			switch reflect.ValueOf(value).Kind() {
			case reflect.Map, reflect.Slice, reflect.Ptr:
				vv := map[interface{}]interface{}{}

				dynamic.Each(value, func(key interface{}, value interface{}) bool {
					if key == "_" {
						dynamic.Each(ReflectValue(app, ctx, value), func(key interface{}, value interface{}) bool {
							vv[key] = value
							return true
						})
					} else {
						vv[key] = ReflectValue(app, ctx, value)
					}
					return true
				})

				return vv
			}
		}

	}
	return value
}
