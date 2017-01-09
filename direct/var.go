package direct

import (
	"strings"
)

type Var struct {
	Direct
}

func (D *Var) Exec(ctx IContext) error {

	options := D.Options()

	v, ok := options["keys"]

	if ok {
		vv, ok := v.(string)
		if ok {
			keys := strings.Split(vv, ".")
			v, ok = options["value"]
			if ok {
				v = ReflectValue(D.App(), ctx, v)
				ctx.Set(keys, v)
			}
		}
	}

	return D.Done(ctx, "done")
}
