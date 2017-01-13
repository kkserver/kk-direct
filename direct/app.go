package direct

import (
	"log"
)

type App struct {
	options Options
	directs map[string]IDirect
}

func NewApp(options Options) IApp {
	v := App{}
	v.options = options
	v.directs = map[string]IDirect{}
	for key, value := range options {
		name, ok := key.(string)
		if ok {
			opt, ok := value.(Options)
			if ok {
				d, err := Open(opt)
				if err != nil {
					log.Println(err)
				} else {
					d.SetApp(&v)
					v.directs[name] = d
				}
			}
		}
	}
	return &v
}

func (A *App) Open(options Options) (IDirect, error) {

	name := options.Name()
	v, ok := A.directs[name]

	if ok {
		return v, nil
	}

	v, err := Open(options)

	if err != nil {
		return nil, err
	}

	v.SetApp(A)

	return v, nil
}

func (A *App) Exec(ctx IContext) error {
	v, ok := A.directs["in"]
	if ok {
		return v.Exec(ctx)
	}
	return NewError(ERROR_UNKNOWN, "Not Found [in] direct")
}

func (A *App) Each(fn func(name string, direct IDirect) bool) {
	for name, direct := range A.directs {
		if !fn(name, direct) {
			break
		}
	}
}
