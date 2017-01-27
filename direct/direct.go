package direct

import (
	"fmt"
	"log"
)

var ErrorKeys = []string{"error"}

const ERROR_UNKNOWN = 0x0ffff

type Error struct {
	Errno  int
	Errmsg string
}

func (E *Error) Error() string {
	return fmt.Sprintf("[%d] %s", E.Errno, E.Errmsg)
}

func NewError(errno int, errmsg string) *Error {
	return &Error{errno, errmsg}
}

type RedirectError struct {
	Url string
}

func (E *RedirectError) Error() string {
	return E.Url
}

func NewRedirectError(url string) *RedirectError {
	return &RedirectError{url}
}

type ContentError struct {
	Content     string
	ContentType string
}

func (E *ContentError) Error() string {
	return E.Content
}

func NewContentError(content string, contentType string) *ContentError {
	return &ContentError{content, contentType}
}

type IContext interface {
	Begin()
	End()
	Get(keys []string) interface{}
	Set(keys []string, value interface{})
}

type IDirect interface {
	App() IApp
	SetApp(app IApp)
	Options() Options
	SetOptions(options Options)
	Exec(ctx IContext) error
	Done(ctx IContext, name string) error
	Fail(ctx IContext, err error) error
	Has(name string) bool
}

type IApp interface {
	Exec(ctx IContext) error
	Open(options Options) (IDirect, error)
	Each(fn func(name string, direct IDirect) bool)
}

type Direct struct {
	app     IApp
	options Options
}

func (D *Direct) App() IApp {
	return D.app
}

func (D *Direct) SetApp(app IApp) {
	D.app = app
}

func (D *Direct) Options() Options {
	if D.options == nil {
		D.options = Options{}
	}
	return D.options
}

func (D *Direct) SetOptions(options Options) {
	D.options = options
}

func (D *Direct) Exec(ctx IContext) error {
	return D.Done(ctx, "done")
}

func (D *Direct) Done(ctx IContext, name string) error {
	log.Println(D.Options().Name(), name)
	v, ok := D.Options()[name]
	if ok {
		vv, ok := v.(Options)
		if ok {
			d, err := D.App().Open(vv)
			if err != nil {
				return err
			}
			return d.Exec(ctx)
		}
	}
	return nil
}

func (D *Direct) Fail(ctx IContext, err error) error {
	log.Println(D.Options().Name(), "fail", err)
	if D.Has("fail") {
		ctx.Set(ErrorKeys, err)
		return D.Done(ctx, "fail")
	}
	return err
}

func (D *Direct) Has(name string) bool {
	v, ok := D.Options()[name]
	if ok {
		_, ok = v.(Options)
		return ok
	}
	return false
}
