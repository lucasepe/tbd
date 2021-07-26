package cmd

import (
	"fmt"

	"github.com/lucasepe/tbd/pkg/data"
	"github.com/lucasepe/tbd/pkg/template"
)

type MarksCmd struct {
	Template string `arg:"positional,required" placeholder:"TEMPLATE"`
}

func (c *MarksCmd) Run() error {
	const maxFileSize int64 = 512 * 1000
	tpl, err := data.Fetch(c.Template, maxFileSize)
	if err != nil {
		return err
	}

	list, _ := template.Marks(string(tpl), "{{", "}}")
	for _, x := range list {
		fmt.Println(x)
	}

	return nil
}
