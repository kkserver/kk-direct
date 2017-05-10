package httpd

import (
	"bufio"
	Y "github.com/go-yaml/yaml"
	"github.com/kkserver/kk-direct/direct"
	KK "github.com/kkserver/kk-direct/direct/kk"
	Lua "github.com/kkserver/kk-direct/direct/lua"
	"github.com/kkserver/kk-direct/direct/view"
	"github.com/kkserver/kk-direct/direct/yaml"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/client"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type App struct {
	app.App
	Client     *client.Service
	Address    string
	Timeout    int
	MaxMemory  int64
	Debug      bool
	WebPath    string
	StaticPath string
	Config     map[string]interface{}

	programs  map[string]direct.IApp
	fs        http.Handler
	static_fs http.Handler
}

func (a *App) Obtain() {
	a.programs = map[string]direct.IApp{}
	a.fs = http.FileServer(http.Dir(a.WebPath))
	a.static_fs = http.FileServer(http.Dir(a.StaticPath))
	log.Println(a.StaticPath)
}

func (a *App) GetProgram(path string) (direct.IApp, error) {

	if a.Debug {

		p, err := yaml.Load(path)

		if err != nil {
			return nil, err
		}

		log.Println(path)

		return p, nil
	}

	p, ok := a.programs[path]

	if ok {
		return p, nil
	}

	p, err := yaml.Load(path)

	if err != nil {
		return nil, err
	}

	a.programs[path] = p

	log.Println(path)

	return p, nil
}

func (a *App) GetContext(w http.ResponseWriter, r *http.Request) direct.IContext {

	ctx := direct.NewContext()

	ctx.Begin()

	ctx.Set(KK.AppKeys, &a)
	ctx.Set(direct.OutputKeys, map[interface{}]interface{}{})
	ctx.Set([]string{"config"}, a.Config)

	Lua.ContextOpenlib(ctx)

	{
		var input interface{} = nil

		if r.Method == "POST" {

			ctype := r.Header.Get("Content-Type")

			if strings.Contains(ctype, "text/json") || strings.Contains(ctype, "application/json") {
				var body = make([]byte, r.ContentLength)
				_, _ = r.Body.Read(body)
				defer r.Body.Close()
				_ = json.Decode(body, &input)
			} else if strings.Contains(ctype, "text/xml") || strings.Contains(ctype, "text/plain") {
				var body = make([]byte, r.ContentLength)
				_, _ = r.Body.Read(body)
				defer r.Body.Close()
				ctx.Set([]string{"content"}, string(body))
			} else if strings.Contains(ctype, "multipart/form-data") {
				r.ParseMultipartForm(a.MaxMemory)
				if r.MultipartForm != nil {
					input = map[interface{}]interface{}{}
					for key, values := range r.MultipartForm.Value {
						dynamic.Set(input, key, values[0])
					}
					for key, values := range r.MultipartForm.File {
						dynamic.Set(input, key, values[0])
					}
				}
			} else {
				r.ParseForm()
				input = map[interface{}]interface{}{}
				for key, values := range r.Form {
					dynamic.Set(input, key, values[0])
				}
			}

		} else {

			r.ParseForm()

			input = map[interface{}]interface{}{}
			for key, values := range r.Form {
				dynamic.Set(input, key, values[0])
			}

		}

		if input != nil {
			ctx.Set([]string{"input"}, input)
		}

	}

	var code = ""

	{
		var ip = r.Header.Get("X-CLIENT-IP")

		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}

		if ip == "" {
			ip = r.RemoteAddr
		}

		ip = strings.Split(ip, ":")[0]

		var cookie, err = r.Cookie("kk")

		if err != nil {
			var v = http.Cookie{}
			v.Name = "kk"
			v.Value = strconv.FormatInt(time.Now().UnixNano(), 10)
			v.Expires = time.Now().Add(24 * 3600 * time.Second)
			v.HttpOnly = true
			v.MaxAge = 24 * 3600
			v.Path = "/"
			http.SetCookie(w, &v)
			cookie = &v
		}

		code = cookie.Value

		var b, _ = json.Encode(map[string]string{"code": code, "ip": ip,
			"User-Agent": r.Header.Get("User-Agent"),
			"Referer":    r.Header.Get("Referer"),
			"Path":       r.RequestURI,
			"Host":       r.Host,
			"Protocol":   r.Proto})

		var m = kk.Message{"MESSAGE", "", "kk.message.http.request", "text/json", b}

		kk.GetDispatchMain().Async(func() {
			task := client.ClientSendMessageTask{}
			task.Message = m
			app.Handle(a, &task)
		})

		ctx.Set([]string{"clientIp"}, ip)
	}

	ctx.Set([]string{"code"}, code)
	ctx.Set([]string{"method"}, r.Method)
	ctx.Set([]string{"path"}, r.URL.Path)
	ctx.Set([]string{"host"}, r.Host)
	ctx.Set([]string{"protocol"}, r.Proto)
	ctx.Set([]string{"referer"}, r.Referer())
	ctx.Set([]string{"userAgent"}, r.UserAgent())
	if strings.HasPrefix(r.Proto, "HTTPS/") {
		ctx.Set([]string{"url"}, "https://"+r.Host+r.RequestURI)
	} else {
		ctx.Set([]string{"url"}, "http://"+r.Host+r.RequestURI)
	}
	ctx.Set([]string{"uri"}, r.RequestURI)

	return ctx
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.HasSuffix(r.URL.Path, ".yaml") || strings.HasSuffix(r.URL.Path, ".yml") || strings.HasSuffix(r.URL.Path, ".htm") {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	} else if strings.HasPrefix(r.URL.Path, "/static/") {
		a.static_fs.ServeHTTP(w, r)
	} else if strings.HasSuffix(r.URL.Path, ".json") {

		path := a.WebPath + r.URL.Path[0:len(r.URL.Path)-5] + ".yaml"

		p, err := a.GetProgram(path)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		} else {
			ctx := a.GetContext(w, r)
			defer Lua.ContextCloselib(ctx)
			defer ctx.End()
			err = p.Exec(ctx)
			if err != nil {
				re, ok := err.(*direct.RedirectError)
				if ok {
					w.Header().Add("Location", re.Url)
					w.WriteHeader(http.StatusTemporaryRedirect)
					w.Write([]byte(""))
				} else {
					ce, ok := err.(*direct.ContentError)
					if ok {
						w.Header().Add("Content-Type", ce.ContentType)
						w.Write([]byte(ce.Content))
					} else {
						ee, ok := err.(*direct.Error)
						if ok {
							b, _ := json.Encode(map[interface{}]interface{}{"errno": ee.Errno, "errmsg": ee.Errmsg})
							w.Header().Add("Content-Type", "application/json; charset=utf-8")
							w.Write(b)
						} else {
							b, _ := json.Encode(map[interface{}]interface{}{"errno": direct.ERROR_UNKNOWN, "errmsg": err.Error()})
							w.Header().Add("Content-Type", "application/json; charset=utf-8")
							w.Write(b)
						}
					}
				}

			} else {

				b, _ := json.Encode(ctx.Get(direct.OutputKeys))
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
				w.Write(b)

			}
		}

	} else if strings.HasSuffix(r.URL.Path, ".lhtml") {

		path := a.WebPath + r.URL.Path[0:len(r.URL.Path)-6] + ".yaml"

		p, err := a.GetProgram(path)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		} else {
			ctx := a.GetContext(w, r)
			defer Lua.ContextCloselib(ctx)
			defer ctx.End()
			err = p.Exec(ctx)
			if err != nil {

				re, ok := err.(*direct.RedirectError)
				if ok {
					w.Header().Add("Location", re.Url)
					w.WriteHeader(http.StatusTemporaryRedirect)
					w.Write([]byte(""))
				} else {
					ce, ok := err.(*direct.ContentError)
					if ok {
						w.Header().Add("Content-Type", ce.ContentType)
						w.Write([]byte(ce.Content))
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(err.Error()))
					}
				}

			} else {

				vv := ctx.Get(view.ViewKeys)

				if vv == nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Not Found"))
				} else {
					v, ok := vv.(*view.View)
					if ok {

						if v.ContentType == "" {
							v.ContentType = "text/html; charset=utf-8"
						}
						w.Header().Add("Content-Type", v.ContentType)
						w.Write(v.Content)
					} else {
						w.WriteHeader(http.StatusNotFound)
						w.Write([]byte("Not Found"))
					}
				}

			}
		}
	} else if strings.HasSuffix(r.URL.Path, "*.doc") {

		root := a.WebPath + r.URL.Path[0:len(r.URL.Path)-5]

		items := []interface{}{}

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if !info.IsDir() {
				if strings.HasSuffix(path, ".yaml") {
					v, _ := filepath.Rel(root, path)
					items = append(items, v[0:len(v)-5]+".doc")
				}
			}

			return nil
		})

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		} else {
			data, _ := json.Encode(map[interface{}]interface{}{"items": items})
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.Write(data)
		}

	} else if strings.HasSuffix(r.URL.Path, ".doc") {

		path := a.WebPath + r.URL.Path[0:len(r.URL.Path)-4] + ".yaml"

		err := func() error {

			fd, err := os.Open(path)

			if err != nil {
				return err
			}

			defer fd.Close()

			rd := bufio.NewReader(fd)

			data, err := rd.ReadBytes(0)

			if err != nil && err != io.EOF {
				return err
			}

			options := direct.Options{}

			err = Y.Unmarshal(data, options)

			if err != nil {
				return err
			}

			data, err = json.Encode(dynamic.Get(options, "doc"))

			if err != nil {
				return err
			}

			return direct.NewContentError(string(data), "application/json; charset=utf-8")

		}()

		v, ok := err.(*direct.ContentError)

		if ok {
			w.Header().Add("Content-Type", v.ContentType)
			w.Write([]byte(v.Content))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		}

	} else if strings.HasSuffix(r.URL.Path, "/") {

		path := a.WebPath + r.URL.Path + "index.yaml"

		p, err := a.GetProgram(path)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		} else {
			ctx := a.GetContext(w, r)
			defer Lua.ContextCloselib(ctx)
			defer ctx.End()
			err = p.Exec(ctx)
			if err != nil {
				re, ok := err.(*direct.RedirectError)
				if ok {
					w.Header().Add("Location", re.Url)
					w.WriteHeader(http.StatusTemporaryRedirect)
					w.Write([]byte(""))
				} else {
					ce, ok := err.(*direct.ContentError)
					if ok {
						w.Header().Add("Content-Type", ce.ContentType)
						w.Write([]byte(ce.Content))
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(err.Error()))
					}
				}
			} else {

				vv := ctx.Get(view.ViewKeys)

				if vv == nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("Not Found"))
				} else {
					v, ok := vv.(*view.View)
					if ok {
						if v.ContentType == "" {
							v.ContentType = "text/html; charset=utf-8"
						}
						w.Header().Add("Content-Type", v.ContentType)
						w.Write(v.Content)
					} else {
						w.WriteHeader(http.StatusNotFound)
						w.Write([]byte("Not Found"))
					}
				}

			}
		}

	} else {
		a.fs.ServeHTTP(w, r)
	}
}
