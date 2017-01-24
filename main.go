package main

import (
	"github.com/kkserver/kk-direct/direct"
	KK "github.com/kkserver/kk-direct/direct/kk"
	Lua "github.com/kkserver/kk-direct/direct/lua"
	"github.com/kkserver/kk-direct/direct/view"
	"github.com/kkserver/kk-direct/direct/yaml"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/client"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type MainApp struct {
	app.App
	Client  *client.Service
	Address string
	Timeout int
	Debug   bool
}

func main() {

	log.SetFlags(log.Llongfile | log.LstdFlags)

	env := "./config/env.ini"

	if len(os.Args) > 1 {
		env = os.Args[1]
	}

	a := MainApp{}

	err := app.Load(&a, "./app.ini")

	if err != nil {
		log.Panicln(err)
	}

	err = app.Load(&a, env)

	if err != nil {
		log.Panicln(err)
	}

	app.Obtain(&a)

	app.Handle(&a, &app.InitTask{})

	yaml.Openlib()
	Lua.Openlib()
	view.Openlib()
	KK.Openlib()
	direct.Openlib()

	go func() {

		programs := map[string]direct.IApp{}

		getProgram := func(path string) (direct.IApp, error) {

			if a.Debug {

				p, err := yaml.Load(path)

				if err != nil {
					return nil, err
				}

				log.Println(path)

				return p, nil
			}

			p, ok := programs[path]

			if ok {
				return p, nil
			}

			p, err := yaml.Load(path)

			if err != nil {
				return nil, err
			}

			programs[path] = p

			log.Println(path)

			return p, nil
		}

		getContext := func(w http.ResponseWriter, r *http.Request) direct.IContext {

			ctx := direct.NewContext()

			ctx.Begin()

			ctx.Set(KK.AppKeys, &a)
			ctx.Set(direct.OutputKeys, map[string]interface{}{})

			Lua.ContextOpenlib(ctx)

			input := map[string]interface{}{}

			if r.Method == "POST" {

				if r.Header.Get("Content-Type") == "text/json" {
					var body = make([]byte, r.ContentLength)
					_, _ = r.Body.Read(body)
					defer r.Body.Close()
					_ = json.Decode(body, input)
				} else {
					r.ParseForm()
					for key, values := range r.Form {
						input[key] = values[0]
					}
				}

			} else {

				r.ParseForm()

				for key, values := range r.Form {
					input[key] = values[0]
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
					app.Handle(&a, &task)
				})

				ctx.Set([]string{"clientIp"}, ip)
			}

			ctx.Set([]string{"code"}, code)
			ctx.Set([]string{"method"}, r.Method)
			ctx.Set([]string{"path"}, r.URL.Path)
			ctx.Set([]string{"host"}, r.Host)
			ctx.Set([]string{"url"}, r.URL)
			ctx.Set([]string{"uri"}, r.RequestURI)
			ctx.Set([]string{"input"}, input)

			return ctx
		}

		fs := http.FileServer(http.Dir("./web"))
		static_fs := http.FileServer(http.Dir("."))

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			if strings.HasSuffix(r.URL.Path, ".yaml") {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Not Found"))
			} else if strings.HasPrefix(r.URL.Path, "/static/") {
				static_fs.ServeHTTP(w, r)
			} else if strings.HasSuffix(r.URL.Path, ".json") {

				path := "./web" + r.URL.Path[0:len(r.URL.Path)-5] + ".yaml"

				p, err := getProgram(path)

				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(err.Error()))
				} else {
					ctx := getContext(w, r)
					defer Lua.ContextCloselib(ctx)
					defer ctx.End()
					err = p.Exec(ctx)
					if err != nil {
						ee, ok := err.(*direct.Error)
						if ok {
							b, _ := json.Encode(map[string]interface{}{"errno": ee.Errno, "errmsg": ee.Errmsg})
							w.Header().Add("Content-Type", "application/json; charset=utf-8")
							w.Write(b)
						} else {
							b, _ := json.Encode(map[string]interface{}{"errno": direct.ERROR_UNKNOWN, "errmsg": err.Error()})
							w.Header().Add("Content-Type", "application/json; charset=utf-8")
							w.Write(b)
						}

					} else {

						b, _ := json.Encode(ctx.Get(direct.OutputKeys))
						w.Header().Add("Content-Type", "application/json; charset=utf-8")
						w.Write(b)

					}
				}

			} else if strings.HasSuffix(r.URL.Path, ".lhtml") {

				path := "./web" + r.URL.Path[0:len(r.URL.Path)-6] + ".yaml"

				p, err := getProgram(path)

				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(err.Error()))
				} else {
					ctx := getContext(w, r)
					defer Lua.ContextCloselib(ctx)
					defer ctx.End()
					err = p.Exec(ctx)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(err.Error()))
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
			} else if strings.HasSuffix(r.URL.Path, "/") {

				path := "./web" + r.URL.Path + "index.yaml"

				p, err := getProgram(path)

				if err != nil {
					fs.ServeHTTP(w, r)
				} else {
					ctx := getContext(w, r)
					defer Lua.ContextCloselib(ctx)
					defer ctx.End()
					err = p.Exec(ctx)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(err.Error()))
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
				fs.ServeHTTP(w, r)
			}

		})

		log.Println("httpd " + a.Address)

		log.Fatal(http.ListenAndServe(a.Address, nil))

	}()

	kk.DispatchMain()

	app.Recycle(&a)

}
