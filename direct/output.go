package direct

import (
	"github.com/kkserver/kk-lib/kk/dynamic"
	"reflect"
	"strings"
)

var OutputKeys = []string{"output"}
var ObjectKeys = []string{"object"}

type Output struct {
	Direct
}

func toObject(output *Output, ctx IContext, fields interface{}, v interface{}) interface{} {

	rs := map[interface{}]interface{}{}

	dynamic.Each(fields, func(key interface{}, opt interface{}) bool {
		s, ok := opt.(string)
		if ok {
			rs[key] = dynamic.Get(v, s)
		} else {
			o, ok := opt.(Options)
			if ok {
				d, err := output.App().Open(o)
				if err == nil {
					ctx.Begin()
					ctx.Set(ObjectKeys, v)
					ctx.Set(OutputKeys, map[interface{}]interface{}{})
					ctx.Set(ResultKeys, Nil)
					d.Exec(ctx)
					vvv := ctx.Get(OutputKeys)
					ctx.End()
					rs[key] = vvv
				}
			}
		}
		return true
	})

	return rs
}

func (D *Output) Exec(ctx IContext) error {

	options := D.Options()

	output := ctx.Get(OutputKeys)

	if output == nil {
		output = map[interface{}]interface{}{}
		ctx.Set(OutputKeys, output)
	}

	v, ok := options["keys"]

	if ok {
		keys, ok := v.(string)
		if ok {
			v, ok = options["value"]
			if ok {

				v = ReflectValue(D.App(), ctx, v)

				{
					fields, ok := options["fields"]

					if ok {
						vv := reflect.ValueOf(v)
						switch vv.Kind() {
						case reflect.Slice:
							{
								rs := []interface{}{}
								for i := 0; i < vv.Len(); i++ {
									vvv := vv.Index(i)
									if vvv.CanInterface() {
										rs = append(rs, toObject(D, ctx, fields, vvv.Interface()))
									}
								}
								v = rs
							}
						case reflect.Ptr, reflect.Map:
							{
								v = toObject(D, ctx, fields, v)
							}
						}
					}

				}

				if keys == "_" {

					dynamic.Each(v, func(key interface{}, value interface{}) bool {
						dynamic.Set(output, dynamic.StringValue(key, ""), value)
						return true
					})

				} else {
					dynamic.SetWithKeys(output, strings.Split(keys, "."), v)
				}

			}
		}
	}

	return D.Done(ctx, "done")
}
