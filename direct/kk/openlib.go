package kk

import (
	"github.com/kkserver/kk-direct/direct"
	"strings"
)

func Openlib() {
	direct.Use(func(name string, options direct.Options) (direct.IDirect, error) {
		if strings.HasPrefix(name, "kk.") {
			v := Direct{}
			v.SetOptions(options)
			return &v, nil
		}
		return nil, nil
	})
}
