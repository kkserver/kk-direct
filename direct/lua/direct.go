package lua

import (
	"fmt"
	"github.com/aarzilli/golua/lua"
	"github.com/kkserver/kk-direct/direct"
)

type Direct struct {
	direct.Direct
}

func (D *Direct) Exec(ctx direct.IContext) error {

	options := D.Options()

	v, ok := options["content"]

	if ok {

		code, ok := v.(string)

		if ok {

			v = ctx.Get(LuaKeys)

			if v != nil {
				L, ok := v.(*lua.State)
				if ok {
					rs := func() interface{} {
						if L.LoadString(fmt.Sprintf("return %s", code)) == 0 {

							err := L.Call(0, 1)

							if err != nil {
								return err.Error()
							}

							if L.IsFunction(-1) {

								err = L.Call(0, 1)

								if err != nil {
									return err.Error()
								}
							}

							var vv interface{} = LuaToValue(L, -1)

							L.Pop(1)

							return vv
						} else {
							vv := L.ToString(-1)
							L.Pop(1)
							return vv
						}
					}()
					ctx.Set(direct.ResultKeys, rs)
				}
			}
		}
	}

	return nil
}

func LuaToValue(L *lua.State, i int) interface{} {

	var vv interface{} = nil

	if L.IsString(i) {
		vv = L.ToString(i)
	} else if L.IsGoStruct(i) {
		vv = L.ToGoStruct(i)
	} else if L.IsNumber(i) {
		vv = L.ToNumber(i)
	} else if L.IsBoolean(i) {
		vv = L.ToBoolean(i)
	} else if L.IsTable(i) {

		L.PushValue(i)

		idx := 0
		size := 0

		m := map[string]interface{}{}
		a := []interface{}{}

		L.PushNil()

		for L.Next(-2) != 0 {

			t := L.Type(-2)

			if t == lua.LUA_TNUMBER {
				if idx == L.ToInteger(-2)-1 {
					a = append(a, LuaToValue(L, -1))
					idx = idx + 1
				}
				m[fmt.Sprintf("%d", L.ToInteger(-2))] = LuaToValue(L, -1)
			} else if t == lua.LUA_TSTRING {
				m[L.ToString(-2)] = LuaToValue(L, -1)
			}

			size = size + 1

			L.Pop(1)

		}

		if idx != 0 && idx == size {
			vv = a
		} else {
			vv = m
		}

		L.Pop(1)
	}

	return vv
}
