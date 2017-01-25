package direct

import (
	"github.com/kkserver/kk-lib/kk/dynamic"
)

type Context struct {
	values []map[interface{}]interface{}
}

func NewContext() IContext {
	return &Context{}
}

func (C *Context) Begin() {
	if C.values == nil {
		C.values = []map[interface{}]interface{}{map[interface{}]interface{}{}}
	} else {
		C.values = append(C.values, map[interface{}]interface{}{})
	}
}

func (C *Context) End() {
	if C.values != nil && len(C.values) > 1 {
		C.values = C.values[0 : len(C.values)-1]
	}
}

func (C *Context) Set(keys []string, value interface{}) {
	if C.values != nil && len(C.values) > 0 {
		vs := C.values[len(C.values)-1]
		dynamic.SetWithKeys(vs, keys, value)
	}

}

func (C *Context) Get(keys []string) interface{} {
	if C.values != nil && len(C.values) > 0 {
		idx := len(C.values) - 1
		ikey := 0
		for ikey < len(keys) && idx > 0 && keys[ikey] == "" {
			idx = idx - 1
			ikey = ikey + 1
		}
		for idx >= 0 {
			vs := C.values[idx]
			v := dynamic.GetWithKeys(vs, keys[ikey:])
			if v != nil {
				return v
			}
			idx = idx - 1
		}
	}
	return nil
}
