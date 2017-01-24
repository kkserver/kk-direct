package lua

import (
	"github.com/aarzilli/golua/lua"
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	"reflect"
)

var LuaKeys = []string{"lua"}

func ContextOpenlib(ctx direct.IContext) {

	L := lua.NewState()
	L.OpenLibs()

	L.PushGoFunction(func(L *lua.State) int {

		keys := []string{}
		top := L.GetTop()

		for i := 0; i < top; i++ {
			keys = append(keys, L.ToString(-top+i))
		}

		vv := ctx.Get(keys)

		LuaPushValue(L, vv)

		return 1
	})

	L.SetGlobal("get")

	L.PushGoFunction(func(L *lua.State) int {

		keys := []string{}
		top := L.GetTop()

		for i := 0; i < top; i++ {
			keys = append(keys, L.ToString(-top+i))
		}

		vv := ctx.Get(keys)

		L.PushString(dynamic.StringValue(vv, ""))

		return 1
	})

	L.SetGlobal("getString")

	L.PushGoFunction(func(L *lua.State) int {

		keys := []string{}
		top := L.GetTop()

		for i := 0; i < top; i++ {
			keys = append(keys, L.ToString(-top+i))
		}

		vv := ctx.Get(keys)

		L.PushInteger(dynamic.IntValue(vv, 0))

		return 1
	})

	L.SetGlobal("getInteger")

	L.PushGoFunction(func(L *lua.State) int {

		keys := []string{}
		top := L.GetTop()

		for i := 0; i < top; i++ {
			keys = append(keys, L.ToString(-top+i))
		}

		vv := ctx.Get(keys)

		L.PushNumber(dynamic.FloatValue(vv, 0))

		return 1
	})

	L.SetGlobal("getNumber")

	L.PushGoFunction(func(L *lua.State) int {

		keys := []string{}
		top := L.GetTop()

		for i := 0; i < top; i++ {
			keys = append(keys, L.ToString(-top+i))
		}

		vv := ctx.Get(keys)

		L.PushBoolean(dynamic.BooleanValue(vv, false))

		return 1
	})

	L.SetGlobal("getBoolean")

	L.NewTable()

	L.PushString("encode")
	L.PushGoFunction(func(L *lua.State) int {

		keys := []string{}
		top := L.GetTop()

		for i := 0; i < top; i++ {
			keys = append(keys, L.ToString(-top+i))
		}

		vv := ctx.Get(keys)

		b, _ := json.Encode(vv)

		L.PushString(string(b))

		return 1
	})
	L.RawSet(-3)

	L.PushString("decode")
	L.PushGoFunction(func(L *lua.State) int {

		keys := []string{}
		top := L.GetTop()

		for i := 0; i < top; i++ {
			keys = append(keys, L.ToString(-top+i))
		}

		vv := ctx.Get(keys)

		var v interface{} = nil

		err := json.Decode([]byte(dynamic.StringValue(vv, "{}")), &v)

		log.Println(vv, keys)

		if err != nil {
			L.PushString(err.Error())
		} else {
			LuaPushValue(L, v)
		}

		return 1
	})
	L.RawSet(-3)

	L.SetGlobal("json")

	ctx.Set(LuaKeys, L)
}

func ContextCloselib(ctx direct.IContext) {
	v := ctx.Get(LuaKeys)
	if v != nil {
		L, ok := v.(*lua.State)
		if ok {
			L.Close()
		}
		ctx.Set(LuaKeys, nil)
	}
}

func Openlib() {
	direct.UseWithType("^lua", reflect.TypeOf(Direct{}))
}
