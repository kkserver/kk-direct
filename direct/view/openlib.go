package view

import (
	"github.com/kkserver/kk-direct/direct"
	"strings"
)

func Openlib() {
	direct.Use(func(name string, options direct.Options) (direct.IDirect, error) {
		if strings.HasSuffix(name, ".html") || strings.HasSuffix(name, ".htm") {
			v := Direct{}
			v.ContentType = "text/html; charset=utf-8"
			v.SetOptions(options)
			return &v, nil
		}
		return nil, nil
	})
}
