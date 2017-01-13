package yaml

import (
	"bufio"
	"fmt"
	Y "github.com/go-yaml/yaml"
	"github.com/kkserver/kk-direct/direct"
	"io"
	"os"
)

func Load(path string) (direct.IApp, error) {

	fd, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer fd.Close()
	rd := bufio.NewReader(fd)

	data, err := rd.ReadBytes(0)

	if err != nil && err != io.EOF {
		return nil, err
	}

	options := direct.Options{}

	err = Y.Unmarshal(data, options)

	if err != nil {
		return nil, direct.NewError(direct.ERROR_UNKNOWN, fmt.Sprintf("%s %s", path, err.Error()))
	}

	return direct.NewApp(options), nil
}
