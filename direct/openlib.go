package direct

import (
	"fmt"
	"reflect"
	"strings"
)

type Library func(name string, options Options) (IDirect, error)

var librarys = []Library{}

func Use(lib Library) {
	librarys = append(librarys, lib)
}

func UseWithType(name string, directType reflect.Type) {
	Use(func(n string, options Options) (IDirect, error) {
		if n == name {
			v, ok := reflect.New(directType).Interface().(IDirect)
			if ok {
				return v, nil
			}
			return nil, NewError(ERROR_UNKNOWN, fmt.Sprintf("%s Not Implement IDirect", directType.Name))
		}
		return nil, nil
	})
}

func Open(options Options) (IDirect, error) {

	name := options.Name()

	for _, lib := range librarys {

		v, err := lib(name, options)

		if err != nil {
			return nil, err
		}

		if v != nil {
			v.SetOptions(options)
			return v, nil
		}
	}

	return nil, NewError(ERROR_UNKNOWN, fmt.Sprintf("Not Open %s", name))
}

func Openlib() {
	UseWithType("^direct", reflect.TypeOf(Direct{}))
	UseWithType("^outlet", reflect.TypeOf(Outlet{}))
	UseWithType("^output", reflect.TypeOf(Output{}))
	UseWithType("^var", reflect.TypeOf(Var{}))
	UseWithType("^required", reflect.TypeOf(Param{}))
	UseWithType("^regexp", reflect.TypeOf(Param{}))
	UseWithType("^int", reflect.TypeOf(Param{}))
	UseWithType("^float", reflect.TypeOf(Param{}))
	UseWithType("^json", reflect.TypeOf(Param{}))
	UseWithType("^date", reflect.TypeOf(Param{}))
	UseWithType("^datetime", reflect.TypeOf(Param{}))
	UseWithType("^day", reflect.TypeOf(Param{}))
	UseWithType("^now", reflect.TypeOf(Param{}))

	Use(func(name string, options Options) (IDirect, error) {

		if strings.HasPrefix(name, "^day") {
			return &Param{}, nil
		}

		if strings.HasPrefix(name, "^datetime") {
			return &Param{}, nil
		}

		if strings.HasPrefix(name, "^date") {
			return &Param{}, nil
		}

		if strings.HasPrefix(name, "^now") {
			return &Param{}, nil
		}

		return nil, nil
	})

	UseWithType("^redirect", reflect.TypeOf(Redirect{}))
	UseWithType("^content", reflect.TypeOf(Content{}))
}
