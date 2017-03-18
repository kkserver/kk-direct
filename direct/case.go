package direct

import (
	"github.com/kkserver/kk-lib/kk/dynamic"
)

type Case struct {
	Direct
}

func (D *Case) Exec(ctx IContext) error {

	options := D.Options()

	when := dynamic.StringValue(ReflectValue(D.App(), ctx, dynamic.Get(options, "when")), "")

	if when != "" && D.Has(when) {
		return D.Done(ctx, when)
	}

	return D.Done(ctx, "done")
}
