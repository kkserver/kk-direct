package direct

import (
	"github.com/kkserver/kk-lib/kk/dynamic"
)

type Content struct {
	Direct
}

func (D *Content) Exec(ctx IContext) error {
	content := ReflectValue(D.App(), ctx, dynamic.Get(D.Options(), "content"))
	ctype := ReflectValue(D.App(), ctx, dynamic.Get(D.Options(), "type"))
	return NewContentError(dynamic.StringValue(content, ""), dynamic.StringValue(ctype, "text/plain"))
}
