package kk

import (
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lib/kk/dynamic"
)

var RouteKeys = []string{"kk", "route"}

type Route struct {
	direct.Direct
}

func (D *Route) Exec(ctx direct.IContext) error {

	route := ctx.Get(RouteKeys)

	if route == nil {
		route = map[interface{}]interface{}{}
		ctx.Set(RouteKeys, route)
	}

	options := D.Options()

	v, ok := options["options"]

	if ok {

		dynamic.Each(v, func(key interface{}, value interface{}) bool {
			dynamic.Set(route, dynamic.StringValue(key, ""), value)
			return true
		})
	}

	return D.Done(ctx, "done")
}
