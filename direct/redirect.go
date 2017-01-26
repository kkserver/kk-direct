package direct

import (
	"github.com/kkserver/kk-lib/kk/dynamic"
)

type Redirect struct {
	Direct
}

func (D *Redirect) Exec(ctx IContext) error {
	url := ReflectValue(D.App(), ctx, dynamic.Get(D.Options(), "url"))
	return NewRedirectError(dynamic.StringValue(url, ""))
}
