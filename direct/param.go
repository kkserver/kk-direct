package direct

import (
	"fmt"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"math"
	"regexp"
	"strings"
	"time"
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
			return D.Fail(ctx, NewError(errno, errmsg))
		}
	case "^regexp":

		pattern, err := regexp.Compile(dynamic.StringValue(dynamic.Get(options, "pattern"), ""))

		if err != nil {
			return err
		}

		if !pattern.MatchString(dynamic.StringValue(vv, "")) {
			return D.Fail(ctx, NewError(errno, errmsg))
		}

	case "^int":
		min := dynamic.IntValue(dynamic.Get(options, "min"), math.MinInt64)
		max := dynamic.IntValue(dynamic.Get(options, "max"), math.MaxInt64)
		vvv := dynamic.IntValue(vv, 0)

		if vvv < min || vvv > max {
			return D.Fail(ctx, NewError(errno, errmsg))
		}

		vv = vvv
	case "^float":
		min := dynamic.FloatValue(dynamic.Get(options, "min"), float64(math.MinInt64))
		max := dynamic.FloatValue(dynamic.Get(options, "max"), math.MaxFloat64)
		vvv := dynamic.FloatValue(vv, 0)

		fmt.Println("^float", vvv, min, max)

		if vvv < min || vvv > max {
			return D.Fail(ctx, NewError(errno, errmsg))
		}

		vv = vvv
	case "^date":
		date, err := time.Parse("2006-01-02", dynamic.StringValue(vv, ""))
		if err != nil {
			return D.Fail(ctx, err)
		}
		vv = date.Unix()
	case "^datetime":
		date, err := time.Parse("2006-01-02 15:04:05", dynamic.StringValue(vv, ""))
		if err != nil {
			return D.Fail(ctx, err)
		}
		vv = date.Unix()
	case "^now":
		vv = time.Now().Unix()
	case "^day":
		now := time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		vv = now.Unix()
	case "^json":
		var data interface{} = nil
		err := json.Decode([]byte(dynamic.StringValue(vv, "{}")), &data)
		if err != nil {
			return D.Fail(ctx, err)
		}
		vv = data
	default:
		if strings.HasPrefix(options.Name(), "^day") {
			now := time.Now()
			now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			now = now.AddDate(0, 0, int(dynamic.IntValue(options.Name()[4:], 0)))
			vv = now.Unix()
		} else if strings.HasPrefix(options.Name(), "^datetime") {
			date, err := time.Parse("2006-01-02 15:04:05", dynamic.StringValue(vv, ""))
			if err != nil {
				return D.Fail(ctx, err)
			}
			vv = date.Unix() + dynamic.IntValue(options.Name()[9:], 0)
		} else if strings.HasPrefix(options.Name(), "^date") {
			date, err := time.Parse("2006-01-02", dynamic.StringValue(vv, ""))
			if err != nil {
				return D.Fail(ctx, err)
			}
			vv = date.Unix() + dynamic.IntValue(options.Name()[5:], 0)
		} else if strings.HasPrefix(options.Name(), "^now") {
			vv = time.Now().Unix() + dynamic.IntValue(options.Name()[4:], 0)
		}
	}

	dynamic.Set(param, key, vv)

	return D.Done(ctx, "done")
}
