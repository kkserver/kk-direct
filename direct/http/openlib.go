package http

import (
	"github.com/kkserver/kk-direct/direct"
	"reflect"
	"strings"
)

func Openlib() {
	direct.Use(func(name string, options direct.Options) (direct.IDirect, error) {
		if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
			v := Direct{}
			v.SetOptions(options)
			return &v, nil
		}
		return nil, nil
	})
}
