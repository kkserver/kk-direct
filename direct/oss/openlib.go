package oss

import (
	"github.com/kkserver/kk-direct/direct"
	"reflect"
)

func Openlib() {
	direct.UseWithType("^oss", reflect.TypeOf(Direct{}))
}
