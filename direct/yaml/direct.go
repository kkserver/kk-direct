package yaml

import (
	"github.com/kkserver/kk-direct/direct"
)

type Direct struct {
	direct.Direct
	app direct.IApp
}

func (D *Direct) Exec(ctx direct.IContext) error {

	if D.app == nil {
		app, err := Load(D.Options().Name())
		if err != nil {
			return err
		}
		app.Each(func(name string, value direct.IDirect) bool {
			v, ok := value.(*direct.Outlet)
			if ok {
				v.On = func(ctx direct.IContext) error {
					return D.Done(ctx, name)
				}
			}
			return true
		})
		D.app = app
	}

	if D.app != nil {
		return D.app.Exec(ctx)
	}

	return nil
}
