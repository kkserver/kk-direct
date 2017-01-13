package direct

import (
	Value "github.com/kkserver/kk-lib/kk/value"
	"reflect"
	"strings"
)

var ParamKeys = []string{"param"}

type Param struct {
	Direct
}

func (D *Output) Exec(ctx IContext) error {

	options := D.Options()

	param := ctx.Get(ParamKeys)

	if param == nil {
		param = map[string]interface{}{}
		ctx.Set(ParamKeys, param)
	}

	v, ok := options["options"]

	if ok {

		opt, ok := v.(direct.Options)

		if ok {
			for key, value := range opt {

				vv, ok := value.(direct.Options)

				if ok {

				}

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

	return D.Done(ctx, "done")
}
