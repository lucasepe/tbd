package cmd

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
)

const (
	description = "A really simple way to create text templates with placeholders."
	banner      = `╔╦╗  ╔╗    ╔╦╗
 ║   ╠╩╗    ║║
 ╩ o ╚═╝ e ═╩╝ efined`
)

type App struct {
	Merge *MergeCmd `arg:"subcommand:merge" help:"combines a template with one or more env files"`
	Marks *MarksCmd `arg:"subcommand:marks" help:"shows all placeholders defined in the specified template"`
	Vars  *VarsCmd  `arg:"subcommand:vars" help:"shows all built-in (and eventually user defined) variables"`
}

func (App) Description() string {
	return fmt.Sprintf("%s\n%s\n", banner, description)
}

func Run() error {
	var app App

	p := arg.MustParse(&app)

	switch {
	case app.Vars != nil:
		return app.Vars.Run()
	case app.Marks != nil:
		return app.Marks.Run()
	case app.Merge != nil:
		return app.Merge.Run()
	default:
		p.WriteHelp(os.Stdout)
	}

	return nil
}
