package direct

type Outlet struct {
	Direct
	On func(ctx IContext) error
}

func (D *Outlet) Exec(ctx IContext) error {
	if D.On != nil {
		return D.On(ctx)
	}
	return nil
}
