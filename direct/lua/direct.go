package lua

import (
	"fmt"
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lua/lua"
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

							if L.Call(0, 1) != 0 {
								err := direct.NewError(direct.ERROR_UNKNOWN, L.ToString(-1))
								L.Pop(1)
								return err
							}

							if L.IsFunction(-1) {

								if L.Call(0, 1) != 0 {
									err := direct.NewError(direct.ERROR_UNKNOWN, L.ToString(-1))
									L.Pop(1)
									return err
								}

							}

							var vv interface{} = L.ToObject(-1)

							L.Pop(1)

							return vv
						} else {
							vv := L.ToString(-1)
							L.Pop(1)
							return vv
						}
					}()
					if rs == nil {
						ctx.Set(direct.ResultKeys, direct.Nil)
					} else {
						ctx.Set(direct.ResultKeys, rs)
					}
				}
			}
		}
	}

	return D.Done(ctx, "done")
}
