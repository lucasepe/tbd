package cmd

import (
	"fmt"
	"sort"

	"github.com/lucasepe/tbd/pkg/table"
)

type VarsCmd struct {
	EnvFiles []string `arg:"positional" placeholder:"ENV_FILE"`
}

func (c *VarsCmd) Run() error {
	meta, err := builtinVars()
	if err != nil {
		return err
	}

	if err := userVars(meta, c.EnvFiles...); err != nil {
		return err
	}

	keys := make([]string, 0, len(meta))
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	tbl := &table.TextTable{}
	tbl.SetHeader("Label", "Value")

	for _, k := range keys {
		tbl.AddRow(k, meta[k])
	}

	fmt.Println(tbl.Draw())

	return nil
}
