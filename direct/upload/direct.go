package upload

import (
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"io"
	"mime/multipart"
	"os"
)

type Direct struct {
	direct.Direct
}

func (D *Direct) Exec(ctx direct.IContext) error {

	options := dynamic.Get(D.Options(), "options")

	v := direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "file"))

	if v != nil {

		path := dynamic.StringValue(direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "path")), "")

		fd, ok := v.(*multipart.FileHeader)

		if ok {

			rd, err := fd.Open()

			if err != nil {
				return D.Fail(ctx, err)
			}

			defer rd.Close()

			wf, err := os.OpenFile(path, os.O_CREATE, 0666)

			if err != nil {
				return D.Fail(ctx, err)
			}

			defer wf.Close()

			data := make([]byte, 20480000)

			for {

				n, err := rd.Read(data)

				if err != nil {
					if err == io.EOF {
						break
					}
					return D.Fail(ctx, err)
				}

				if n == 0 {
					break
				}

				_, err = wf.Write(data[0:n])

				if err != nil {
					return D.Fail(ctx, err)
				}

			}

		} else {
			return D.Fail(ctx, direct.NewError(direct.ERROR_UNKNOWN, "Not multipart/form-data"))
		}
	} else {
		return D.Fail(ctx, direct.NewError(direct.ERROR_UNKNOWN, "Not found file"))
	}

	return D.Done(ctx, "done")
}
