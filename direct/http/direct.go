package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	xhttp "net/http"
	xurl "net/url"
	"strings"
	"time"
)

var ResultKeys = []string{"result"}

var ca *x509.CertPool = nil

type Direct struct {
	direct.Direct
}

func (D *Direct) Exec(ctx direct.IContext) error {

	if ca == nil {
		ca = x509.NewCertPool()
		ca.AppendCertsFromPEM(pemCerts)
	}

	client := &xhttp.Client{
		Transport: &xhttp.Transport{
			TLSClientConfig: &tls.Config{RootCAs: ca},
		},
	}

	options := D.Options()

	url := options.Name()
	method := dynamic.StringValue(dynamic.Get(options, "method"), "GET")
	stype := dynamic.StringValue(dynamic.Get(options, "type"), "application/x-www-form-urlencoded")
	responseType := dynamic.StringValue(dynamic.Get(options, "responseType"), "json")

	data := direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "options"))

	var resp *xhttp.Response
	var err error

	if method == "POST" {

		var body []byte = nil

		if stype == "text/json" || stype == "application/json" {

			body, err = json.Encode(data)

			if err != nil {
				return D.Fail(ctx, err)
			}

		} else {

			idx := 0
			b := bytes.NewBuffer(nil)

			dynamic.Each(data, func(key interface{}, value interface{}) bool {

				if idx != 0 {
					b.WriteString("&")
				}

				b.WriteString(dynamic.StringValue(key, ""))
				b.WriteString("=")
				b.WriteString(xurl.QueryEscape(dynamic.StringValue(value, "")))

				idx = idx + 1

				return true
			})

			body = b.Bytes()
		}

		resp, err = client.Post(url, stype+"; charset=utf-8", body)

	} else {

		idx := 0

		b := bytes.NewBuffer(nil)

		dynamic.Each(data, func(key interface{}, value interface{}) bool {

			if idx != 0 {
				b.WriteString("&")
			}

			b.WriteString(dynamic.StringValue(key, ""))
			b.WriteString("=")
			b.WriteString(xurl.QueryEscape(dynamic.StringValue(value, "")))

			idx = idx + 1

			return true
		})

		idx = strings.Index(url, "?")

		if idx >= 0 {
			if idx+1 == len(url) {
				url = url + b.String()
			} else {
				url = url + "&" + b.String()
			}
		} else {
			url = url + "?" + b.String()
		}

		resp, err = client.Get(url)
	}

	if err != nil {
		return D.Fail(ctx, err)
	}

	if resp.StatusCode == 200 {
		var body = make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(body)
		defer resp.Body.Close()

		if responseType == "json" {
			data := map[interface{}]interface{}{}
			err := json.Decode(body, data)
			if err != nil {
				return D.Fail(ctx, err)
			}
			ctx.Set(ResultKeys, data)
		} else {
			ctx.Set(ResultKeys, string(body))
		}

	} else {
		var body = make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(body)
		defer resp.Body.Close()
		return D.Fail(ctx, errors.New(fmt.Sprintf("[%d] %s", resp.StatusCode, string(body))))
	}

	return D.Done(ctx, "done")
}
