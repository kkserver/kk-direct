package oss

import (
	OSS "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/kkserver/kk-direct/direct"
	"github.com/kkserver/kk-lib/kk/dynamic"
	"log"
	"mime/multipart"
)

type Direct struct {
	direct.Direct
}

func (D *Direct) Exec(ctx direct.IContext) error {

	options := dynamic.Get(D.Options(), "options")

	v := direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "file"))

	if v != nil {

		path := dynamic.StringValue(direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "path")), "")

		client, err := OSS.New(dynamic.StringValue(direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "endpoint")), ""),
			dynamic.StringValue(direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "accessKeyId")), ""),
			dynamic.StringValue(direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "accessKeySecret")), ""))

		if err != nil {
			return D.Fail(ctx, err)
		}

		bucket, err := client.Bucket(dynamic.StringValue(direct.ReflectValue(D.App(), ctx, dynamic.Get(options, "bucket")), ""))

		if err != nil {
			return D.Fail(ctx, err)
		}

		fd, ok := v.(*multipart.FileHeader)

		if ok {

			rd, err := fd.Open()

			if err != nil {
				return D.Fail(ctx, err)
			}

			err = bucket.PutObject(path, rd)

			log.Println("oss", path, err)

			if err != nil {
				return D.Fail(ctx, err)
			}

		} else {
			return D.Fail(ctx, direct.NewError(direct.ERROR_UNKNOWN, "Not multipart/form-data"))
		}
	} else {
		return D.Fail(ctx, direct.NewError(direct.ERROR_UNKNOWN, "Not found file"))
	}

	return D.Done(ctx, "done")
}
