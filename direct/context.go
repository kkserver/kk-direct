package direct

import (
	Value "github.com/kkserver/kk-lib/kk/value"
	"reflect"
)

type Context struct {
	values []map[string]interface{}
}

func NewContext() IContext {
	return &Context{}
}

func (C *Context) Begin() {
	if C.values == nil {
		C.values = []map[string]interface{}{map[string]interface{}{}}
	} else {
		C.values = append(C.values, map[string]interface{}{})
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
		Value.SetWithKeys(reflect.ValueOf(vs), keys, reflect.ValueOf(value))
	}

}

func (C *Context) Get(keys []string) interface{} {
	if C.values != nil && len(C.values) > 0 {
		idx := len(C.values) - 1
		for idx >= 0 {
			vs := C.values[idx]
			v := Value.GetWithKeys(reflect.ValueOf(vs), keys)
			if v.IsValid() && v.CanInterface() && !v.IsNil() {
				return v.Interface()
			}
			idx = idx - 1
		}
	}
	return nil
}
