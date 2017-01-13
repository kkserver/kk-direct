package direct

import (
	"fmt"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"math"
	"regexp"
)

var ParamKeys = []string{"param"}

type Param struct {
	Direct
}

func (D *Param) Exec(ctx IContext) error {

	options := D.Options()

	param := ctx.Get(ParamKeys)

	if param == nil {
		param = map[string]interface{}{}
		ctx.Set(ParamKeys, param)
	}

	key := dynamic.StringValue(dynamic.Get(options, "key"), "")
	vv := ReflectValue(D.App(), ctx, dynamic.Get(options, "value"))
	errno := int(dynamic.IntValue(dynamic.Get(options, "errno"), ERROR_UNKNOWN))
	errmsg := dynamic.StringValue(dynamic.Get(options, "errmsg"), fmt.Sprintf("Param %s fail", key))

	switch options.Name() {
	case "^required":
		if dynamic.IsEmpty(vv) {
			return NewError(errno, errmsg)
		}
	case "^regexp":

		pattern, err := regexp.Compile(dynamic.StringValue(dynamic.Get(options, "pattern"), ""))

		if err != nil {
			return err
		}

		if !pattern.MatchString(dynamic.StringValue(vv, "")) {
			return NewError(errno, errmsg)
		}

	case "^int":
		min := dynamic.IntValue(dynamic.Get(options, "min"), math.MinInt64)
		max := dynamic.IntValue(dynamic.Get(options, "max"), math.MaxInt64)
		vvv := dynamic.IntValue(vv, 0)

		if vvv < min || vvv > max {
			return NewError(errno, errmsg)
		}

		vv = vvv
	case "^float":
		min := dynamic.FloatValue(dynamic.Get(options, "min"), float64(math.MinInt64))
		max := dynamic.FloatValue(dynamic.Get(options, "max"), math.MaxFloat64)
		vvv := dynamic.FloatValue(vv, 0)

		if vvv < min || vvv > max {
			return NewError(errno, errmsg)
		}

		vv = vvv
	case "^json":
		var data interface{} = nil
		err := json.Decode([]byte(dynamic.StringValue(vv, "{}")), &data)
		if err != nil {
			return err
		}
		vv = data
	}

	dynamic.Set(param, key, vv)

	return D.Done(ctx, "done")
}
