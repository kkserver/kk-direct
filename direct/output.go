package direct

import (
	Value "github.com/kkserver/kk-lib/kk/value"
	"reflect"
	"strings"
)

var OutputKeys = []string{"output"}

type Output struct {
	Direct
}

func (D *Output) Exec(ctx IContext) error {

	options := D.Options()

	output := ctx.Get(OutputKeys)

	if output == nil {
		output = map[string]interface{}{}
		ctx.Set(OutputKeys, output)
	}

	v, ok := options["keys"]

	if ok {
		vv, ok := v.(string)
		if ok {
			keys := strings.Split(vv, ".")
			v, ok = options["value"]
			if ok {

				v = ReflectValue(D.App(), ctx, v)

				Value.SetWithKeys(reflect.ValueOf(output), keys, reflect.ValueOf(v))

			}
		}
	}

	return D.Done(ctx, "done")
}
