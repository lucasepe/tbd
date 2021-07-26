package cmd

import (
	"bytes"
	"os"
	"runtime"
	"time"

	"github.com/lucasepe/tbd/pkg/data"
	"github.com/lucasepe/tbd/pkg/dotenv"
	"github.com/lucasepe/tbd/pkg/vcs"
)

const (
	TimeStamp = "TIMESTAMP"
	OS        = "OS"
	ARCH      = "ARCH"
)

func builtinVars() (map[string]string, error) {
	meta := map[string]string{}
	meta[TimeStamp] = time.Now().Local().UTC().Format(time.RFC3339)
	meta[OS] = runtime.GOOS
	meta[ARCH] = runtime.GOARCH

	if cwd, err := os.Getwd(); err == nil {
		vcs.GitRepoMetadata(cwd, meta)
	}

	return meta, nil
}

func userVars(vars map[string]string, envfile ...string) error {
	if len(envfile) <= 0 {
		return nil
	}

	for _, el := range envfile {
		const maxFileSize int64 = 512 * 1000
		buf, err := data.Fetch(el, maxFileSize)
		if err != nil {
			return err
		}

		if err := dotenv.ParseInto(bytes.NewBuffer(buf), vars); err != nil {
			return err
		}
	}

	return nil
}
