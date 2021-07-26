package cmd

import (
	"os"

	"github.com/lucasepe/tbd/pkg/data"
	"github.com/lucasepe/tbd/pkg/template"
)

type MergeCmd struct {
	Template string   `arg:"positional,required" placeholder:"TEMPLATE"`
	EnvFiles []string `arg:"positional" placeholder:"ENV_FILE"`
}

func (c *MergeCmd) Run() error {
	meta, err := builtinVars()
	if err != nil {
		return err
	}

	if err := userVars(meta, c.EnvFiles...); err != nil {
		return err
	}

	const maxFileSize int64 = 512 * 1000
	tpl, err := data.Fetch(c.Template, maxFileSize)
	if err != nil {
		return err
	}

	env := make(map[string]interface{})
	for k, v := range meta {
		env[k] = v
	}

	_, err = template.ExecuteStd(string(tpl), "{{", "}}", os.Stdout, env)
	return err
}
