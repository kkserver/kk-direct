package direct

import (
	"strings"
)

var ResultKeys = []string{"result"}

func ReflectValue(app IApp, ctx IContext, value interface{}) interface{} {

	if value != nil {

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
	return value
}
