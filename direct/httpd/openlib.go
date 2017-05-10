package httpd

import (
	"github.com/kkserver/kk-direct/direct"
	KK "github.com/kkserver/kk-direct/direct/kk"
	Lua "github.com/kkserver/kk-direct/direct/lua"
	"github.com/kkserver/kk-direct/direct/oss"
	"github.com/kkserver/kk-direct/direct/view"
	"github.com/kkserver/kk-direct/direct/yaml"
)

func Openlib() {
	yaml.Openlib()
	Lua.Openlib()
	view.Openlib()
	KK.Openlib()
	oss.Openlib()
	direct.Openlib()
}
