package direct

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"math"
	"regexp"
	"strings"
	"time"
)

var ParamKeys = []string{"param"}
var KeyKeys = []string{"key"}

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
		min := dynamic.IntValue(ReflectValue(D.App(), ctx, dynamic.Get(options, "min")), math.MinInt64)
		max := dynamic.IntValue(ReflectValue(D.App(), ctx, dynamic.Get(options, "max")), math.MaxInt64)
		vvv := dynamic.IntValue(vv, 0)

		if vvv < min || vvv > max {
			return D.Fail(ctx, NewError(errno, errmsg))
		}

		vv = vvv
	case "^float":
		min := dynamic.FloatValue(ReflectValue(D.App(), ctx, dynamic.Get(options, "min")), float64(math.MinInt64))
		max := dynamic.FloatValue(ReflectValue(D.App(), ctx, dynamic.Get(options, "max")), math.MaxFloat64)
		vvv := dynamic.FloatValue(vv, 0)

		fmt.Println("^float", vvv, min, max)

		if vvv < min || vvv > max {
			return D.Fail(ctx, NewError(errno, errmsg))
		}

		vv = vvv
	case "^date":
		date, err := time.ParseInLocation("2006-01-02", dynamic.StringValue(vv, ""), time.Local)
		if err != nil {
			return D.Fail(ctx, err)
		}
		vv = date.Unix()
	case "^datetime":
		date, err := time.ParseInLocation("2006-01-02 15:04:05", dynamic.StringValue(vv, ""), time.Local)
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
	case "^week":
		now := time.Now()
		for now.Weekday() != 0 {
			now = now.AddDate(0, 0, -1)
		}
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		vv = now.Unix()
	case "^json":
		var data interface{} = nil
		err := json.Decode([]byte(dynamic.StringValue(vv, "{}")), &data)
		if err != nil {
			return D.Fail(ctx, err)
		}
		vv = data
	case "^jsonString":
		b, err := json.Encode(vv)
		if err != nil {
			return D.Fail(ctx, err)
		}
		vv = string(b)
	case "^joinString":

		field := dynamic.StringValue(dynamic.Get(options, "field"), "_")
		sep := dynamic.StringValue(dynamic.Get(options, "sep"), ",")
		b := bytes.NewBuffer(nil)
		idx := 0

		dynamic.Each(vv, func(key interface{}, value interface{}) bool {

			if idx != 0 {
				b.WriteString(sep)
			}

			if field == "_" {
				b.WriteString(dynamic.StringValue(value, ""))
			} else {
				b.WriteString(dynamic.StringValue(dynamic.Get(value, field), ""))
			}

			idx = idx + 1

			return true
		})

		vv = b.String()
	case "^array":

		vs := []interface{}{}

		var each IDirect = nil
		var err error = nil

		{
			vvv := dynamic.Get(options, "each")
			if vvv != nil {
				o, ok := vvv.(Options)
				if ok {
					each, err = D.App().Open(o)
					if err != nil {
						return D.Fail(ctx, err)
					}
				}
			}
		}

		dynamic.Each(vv, func(key interface{}, value interface{}) bool {
			if each == nil {
				vs = append(vs, value)
			} else {
				ctx.Begin()
				ctx.Set(ObjectKeys, value)
				ctx.Set(KeyKeys, key)
				ctx.Set(OutputKeys, map[interface{}]interface{}{})
				ctx.Set(ParamKeys, map[interface{}]interface{}{})
				ctx.Set(ResultKeys, Nil)
				err = each.Exec(ctx)
				vvv := ctx.Get(OutputKeys)
				ctx.End()
				if err != nil {
					return false
				}
				vs = append(vs, vvv)
			}
			return true
		})

		if err != nil {
			return D.Fail(ctx, err)
		}

		vv = vs

	case "^object":

		vs := map[interface{}]interface{}{}

		var each IDirect = nil
		var err error = nil

		{
			vvv := dynamic.Get(options, "each")
			if vvv != nil {
				o, ok := vvv.(Options)
				if ok {
					each, err = D.App().Open(o)
					if err != nil {
						return D.Fail(ctx, err)
					}
				}
			}
		}

		dynamic.Each(vv, func(key interface{}, value interface{}) bool {
			if each == nil {
				vs[key] = value
			} else {
				ctx.Begin()
				ctx.Set(ObjectKeys, value)
				ctx.Set(KeyKeys, key)
				ctx.Set(OutputKeys, map[interface{}]interface{}{})
				ctx.Set(ParamKeys, map[interface{}]interface{}{})
				ctx.Set(ResultKeys, Nil)
				err = each.Exec(ctx)
				vvv := ctx.Get(OutputKeys)
				ctx.End()
				if err != nil {
					return false
				}
				vs[key] = vvv
			}
			return true
		})

		if err != nil {
			return D.Fail(ctx, err)
		}

		vv = vs
	case "^uuid":
		vv = kk.UUID()
	default:
		if strings.HasPrefix(options.Name(), "^day") {
			now := time.Now()
			now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			now = now.AddDate(0, 0, int(dynamic.IntValue(options.Name()[4:], 0)))
			vv = now.Unix()
		} else if strings.HasPrefix(options.Name(), "^datetime") {
			date, err := time.ParseInLocation("2006-01-02 15:04:05", dynamic.StringValue(vv, ""), time.Local)
			if err != nil {
				return D.Fail(ctx, err)
			}
			vv = date.Unix() + dynamic.IntValue(options.Name()[9:], 0)
		} else if strings.HasPrefix(options.Name(), "^date") {
			date, err := time.ParseInLocation("2006-01-02", dynamic.StringValue(vv, ""), time.Local)
			if err != nil {
				return D.Fail(ctx, err)
			}
			vv = date.Unix() + dynamic.IntValue(options.Name()[5:], 0)
		} else if strings.HasPrefix(options.Name(), "^now") {
			vv = time.Now().Unix() + dynamic.IntValue(options.Name()[4:], 0)
		} else if strings.HasPrefix(options.Name(), "^week") {
			now := time.Now()
			for now.Weekday() != 0 {
				now = now.AddDate(0, 0, -1)
			}
			now = now.AddDate(0, 0, int(dynamic.IntValue(options.Name()[5:], 0)))
			now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			vv = now.Unix()
		}
	}

	dynamic.Set(param, key, vv)

	return D.Done(ctx, "done")
}
