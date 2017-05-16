package upload

import (
	"github.com/kkserver/kk-direct/direct"
	"reflect"
)

func Openlib() {
	direct.UseWithType("^upload", reflect.TypeOf(Direct{}))
}
