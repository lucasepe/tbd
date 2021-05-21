package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lucasepe/tbd/internal/data"
	"github.com/lucasepe/tbd/internal/dotenv"
	"github.com/lucasepe/tbd/internal/template"
)

const (
	maxFileSize int64 = 512 * 1000
	banner            = `╔╦╗ ╔╗  ╔╦╗
 ║  ╠╩╗  ║║
 ╩  ╚═╝ ═╩╝ to be defined`
)

var (
	optVars    string
	optVersion bool

	commit string
)

func main() {
	configureFlags()

	if optVersion {
		fmt.Printf("%s version: %s\n", appName(), commit)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.CommandLine.Usage()
		os.Exit(0)
	}

	tpl, err := data.Fetch(flag.Args()[0], maxFileSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	if len(optVars) == 0 {
		list, _ := template.Marks(string(tpl), "{{", "}}")
		for _, x := range list {
			fmt.Println(x)
		}
		os.Exit(0)
	}

	buf, err := data.Fetch(optVars, maxFileSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	env, err := dotenv.Parse(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	_, err = template.ExecuteStd(string(tpl), "{{", "}}", os.Stdout, env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

func configureFlags() {
	name := appName()

	flag.CommandLine.Usage = func() {
		fmt.Printf("%s\n\n", banner)
		fmt.Print("A really simple way to create text templates with placeholders.\n\n")

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [flags] <TEXT TEMPLATE URI>\n\n", name)

		fmt.Print("EXAMPLES:\n\n")
		fmt.Print(" > List all defined placeholders:\n\n")
		fmt.Printf("     %s [http or local]/path/to/text/template\n\n", name)
		fmt.Print(" > Replace all placeholders with the values ​​defined in the specified file:\n\n")
		fmt.Printf("     %s -vars [http or local]/path/to/variables [http or local]/path/to/text/template\n\n", name)

		fmt.Print("FLAGS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
		fmt.Print("  -help\n\tprints this message\n")
		fmt.Println()

		fmt.Println("crafted with passion by Luca Sepe <luca.sepe@gmail.com>")
	}

	flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)

	flag.CommandLine.StringVar(&optVars, "vars", "", "the file containing the text template placeholders values")
	flag.CommandLine.BoolVar(&optVersion, "v", false, "print current version and exit")

	flag.CommandLine.Parse(os.Args[1:])
}

func appName() string {
	return filepath.Base(os.Args[0])
}
