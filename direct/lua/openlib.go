package lua

import (
	"database/sql"
	"github.com/aarzilli/golua/lua"
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	"reflect"
	"strings"
)

var dbs = map[string]*sql.DB{}

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

	L.NewTable()

	L.PushString("query")
	L.PushGoFunction(func(L *lua.State) int {

		top := L.GetTop()

		if top > 1 {

			name := L.ToString(-top)
			s := L.ToString(-top + 1)
			args := []interface{}{}

			db, ok := dbs[name]

			var err error

			if !ok {

				idx := strings.Index(name, "://")

				if idx > 0 {

					db, err = sql.Open(name[:idx], name[idx+3:])

					if err != nil {
						log.Println("LUA db.query", "fail", err)
						L.PushGoStruct([]interface{}{})
						return 1
					} else {
						db.SetMaxIdleConns(6)
						db.SetMaxOpenConns(20)
						dbs[name] = db
					}

				} else {
					log.Println("LUA db.query", "fail", name)
					L.PushGoStruct([]interface{}{})
					return 1
				}

			}

			for i := 2; i < top; i++ {
				args = append(args, LuaToValue(L, -top+i))
			}

			log.Println("SQL", s, args)

			rows, err := db.Query(s, args...)

			if err != nil {
				log.Println("LUA db.query", "fail", err)
				log.Println("SQL", s, args)
				L.PushGoStruct([]interface{}{})
				return 1
			}

			defer rows.Close()

			var columns []string = nil
			var rs = []interface{}{}

			columns, err = rows.Columns()

			if err != nil {
				log.Println("LUA db.query", "fail", err)
				L.PushGoStruct([]interface{}{})
				return 1
			}

			var values = make([]sql.NullString, len(columns))
			var refs = make([]interface{}, len(columns))

			for i, _ := range columns {
				refs[i] = &values[i]
			}

			for rows.Next() {

				err = rows.Scan(refs...)

				if err != nil {
					log.Println("LUA db.query", "fail", err)
					L.PushGoStruct([]interface{}{})
					return 1
				}

				log.Println(values)

				row := map[interface{}]interface{}{}

				for i, name := range columns {
					vv := values[i]
					if vv.Valid {
						row[name] = vv.String
					}
				}

				rs = append(rs, row)
			}

			L.PushGoStruct(rs)

			return 1
		}

		L.PushGoStruct([]interface{}{})

		return 1
	})

	L.RawSet(-3)

	L.SetGlobal("db")

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
