package view

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/aarzilli/golua/lua"
	"github.com/kkserver/kk-direct/direct"
	Lua "github.com/kkserver/kk-direct/direct/lua"
	Value "github.com/kkserver/kk-lib/kk/value"
	"io"
	"os"
	"reflect"
	"regexp"
)

var ViewKeys = []string{"view"}

type View struct {
	Content     []byte
	ContentType string
}

var viewLogicCodeRegexp, _ = regexp.Compile("\\{\\#.*?\\#\\}")
var viewLogicIncludeRegexp, _ = regexp.Compile("\\<\\!\\-\\-\\ include\\(.*?\\)\\ \\-\\-\\>")

func GetFileContent(path string) (string, error) {

	fd, err := os.Open(path)

	if err != nil {
		return "", err
	}

	rd := bufio.NewReader(fd)

	v, err := rd.ReadString(0)

	fd.Close()

	if err != nil && err != io.EOF {
		return "", err
	}

	data := bytes.NewBuffer(nil)
	i := 0

	for i < len(v) {

		vs := viewLogicIncludeRegexp.FindStringIndex(v[i:])

		if vs != nil {

			if vs[0] > 0 {
				data.WriteString(v[i : i+vs[0]])
			}

			vv, err := GetFileContent(v[i+vs[0]+13 : i+vs[1]-5])

			if err != nil {
				return "", err
			} else {
				data.WriteString(vv)
			}

			i = i + vs[1]

		} else {
			data.WriteString(v[i:])
			break
		}
	}

	return data.String(), nil
}

type Direct struct {
	direct.Direct
	ContentType string

	hasContent bool
	content    string
}

func (D *Direct) Exec(ctx direct.IContext) error {

	options := D.Options()

	if !D.hasContent {

		v, err := GetFileContent(options.Name())

		if err != nil {
			D.content = err.Error()
		} else {
			D.content = v
		}

		D.hasContent = true
	}

	data := bytes.NewBuffer(nil)

	i := 0

	for i < len(D.content) {

		vs := viewLogicCodeRegexp.FindStringIndex(D.content[i:])

		if vs != nil {

			if vs[0] > 0 {
				data.WriteString(D.content[i : i+vs[0]])
			}

			data.WriteString(D.ExecCode(ctx, D.content[i+vs[0]+2:i+vs[1]-2]))

			i = i + vs[1]

		} else {
			data.WriteString(D.content[i:])
			break
		}
	}

	ctx.Set(ViewKeys, &View{data.Bytes(), D.ContentType})

	return nil
}

func (D *Direct) ExecCode(ctx direct.IContext, code string) string {

	v := ctx.Get(Lua.LuaKeys)

	if v != nil {

		L, ok := v.(*lua.State)

		if ok {

			var vv interface{} = nil

			if L.LoadString(fmt.Sprintf("return %s", code)) == 0 {

				err := L.Call(0, 1)

				if err != nil {
					vv = err.Error()
				} else {

					if L.IsFunction(-1) {

						err = L.Call(0, 1)

						if err != nil {
							vv = err.Error()
						} else {

							vv = Lua.LuaToValue(L, -1)

							L.Pop(1)
						}

					} else {

						vv = Lua.LuaToValue(L, -1)

						L.Pop(1)
					}

				}

			} else {
				vv = L.ToString(-1)
				L.Pop(1)
			}

			return Value.StringValue(reflect.ValueOf(vv), "")
		}

	}

	return code
}
