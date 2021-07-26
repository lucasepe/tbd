package template

import (
	"fmt"
	"io"
	"log"
	"net/url"
)

func ExampleExecuteString() {
	// Substitution map.
	// Since "baz" tag is missing in the map, it will be substituted
	// by an empty string.
	m := map[string]interface{}{
		"host": "google.com",     // string - convenient
		"bar":  []byte("foobar"), // byte slice - the fastest

		// TagFunc - flexible value. TagFunc is called only if the given
		// tag exists in the template.
		"query": TagFunc(func(w io.Writer, tag string) (int, error) {
			return w.Write([]byte(url.QueryEscape(tag + "=world")))
		}),
	}

	template := "http://{{host}}/?foo={{bar}}{{bar}}&q={{query}}&baz={{baz}}"

	s, err := ExecuteString(template, "{{", "}}", m)
	if err != nil {
		log.Fatalf("unexpected error when executing template: %s", err)
	}

	fmt.Printf("%s", s)

	// Output:
	// http://google.com/?foo=foobarfoobar&q=query%3Dworld&baz=
}

func ExampleTagFunc() {
	bazSlice := [][]byte{[]byte("123"), []byte("456"), []byte("789")}
	m := map[string]interface{}{
		// Always wrap the function into TagFunc.
		//
		// "baz" tag function writes bazSlice contents into w.
		"baz": TagFunc(func(w io.Writer, tag string) (int, error) {
			var nn int
			for _, x := range bazSlice {
				n, err := w.Write(x)
				if err != nil {
					return nn, err
				}
				nn += n
			}
			return nn, nil
		}),
	}

	template := "foo[baz]bar"

	s, err := ExecuteString(template, "[", "]", m)
	if err != nil {
		log.Fatalf("unexpected error when executing template: %s", err)
	}
	fmt.Printf("%s", s)

	// Output:
	// foo123456789bar
}

func ExampleExecuteFuncString() {
	template := "Hello, [user]! You won [prize]!!! [foobar]"

	s, err := ExecuteFuncString(template, "[", "]", func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "user":
			return w.Write([]byte("John"))
		case "prize":
			return w.Write([]byte("$100500"))
		default:
			return w.Write([]byte(fmt.Sprintf("[unknown tag %q]", tag)))
		}
	})
	if err != nil {
		log.Fatalf("unexpected error when executing template: %s", err)
	}

	fmt.Printf("%s", s)

	// Output:
	// Hello, John! You won $100500!!! [unknown tag "foobar"]
}
