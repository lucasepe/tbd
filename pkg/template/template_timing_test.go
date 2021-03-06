package template

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"text/template"
)

var (
	source             = "http://{{uid}}.foo.bar.com/?cb={{cb}}{{width}}&width={{width}}&height={{height}}&timeout={{timeout}}&uid={{uid}}&subid={{subid}}&ref={{ref}}&empty={{empty}}"
	result             = "http://aaasdf.foo.bar.com/?cb=12341232&width=1232&height=123&timeout=123123&uid=aaasdf&subid=asdfds&ref=http://google.com/aaa/bbb/ccc&empty="
	resultEscaped      = "http://aaasdf.foo.bar.com/?cb=12341232&width=1232&height=123&timeout=123123&uid=aaasdf&subid=asdfds&ref=http%3A%2F%2Fgoogle.com%2Faaa%2Fbbb%2Fccc&empty="
	resultStd          = "http://aaasdf.foo.bar.com/?cb=12341232&width=1232&height=123&timeout=123123&uid=aaasdf&subid=asdfds&ref=http://google.com/aaa/bbb/ccc&empty={{empty}}"
	resultTextTemplate = "http://aaasdf.foo.bar.com/?cb=12341232&width=1232&height=123&timeout=123123&uid=aaasdf&subid=asdfds&ref=http://google.com/aaa/bbb/ccc&empty=<no value>"

	resultBytes             = []byte(result)
	resultEscapedBytes      = []byte(resultEscaped)
	resultStdBytes          = []byte(resultStd)
	resultTextTemplateBytes = []byte(resultTextTemplate)

	m = map[string]interface{}{
		"cb":      []byte("1234"),
		"width":   []byte("1232"),
		"height":  []byte("123"),
		"timeout": []byte("123123"),
		"uid":     []byte("aaasdf"),
		"subid":   []byte("asdfds"),
		"ref":     []byte("http://google.com/aaa/bbb/ccc"),
	}
)

func map2slice(m map[string]interface{}) []string {
	var a []string
	for k, v := range m {
		a = append(a, "{{"+k+"}}", string(v.([]byte)))
	}
	return a
}

func BenchmarkFmtFprintf(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			fmt.Fprintf(&w,
				"http://%[5]s.foo.bar.com/?cb=%[1]s%[2]s&width=%[2]s&height=%[3]s&timeout=%[4]s&uid=%[5]s&subid=%[6]s&ref=%[7]s&empty=",
				m["cb"], m["width"], m["height"], m["timeout"], m["uid"], m["subid"], m["ref"])
			x := w.Bytes()
			if !bytes.Equal(x, resultBytes) {
				b.Fatalf("Unexpected result\n%q\nExpected\n%q\n", x, result)
			}
			w.Reset()
		}
	})
}

func BenchmarkStringsReplace(b *testing.B) {
	mSlice := map2slice(m)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := source
			for i := 0; i < len(mSlice); i += 2 {
				x = strings.Replace(x, mSlice[i], mSlice[i+1], -1)
			}
			if x != resultStd {
				b.Fatalf("Unexpected result\n%q\nExpected\n%q\n", x, resultStd)
			}
		}
	})
}

func BenchmarkStringsReplacer(b *testing.B) {
	mSlice := map2slice(m)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r := strings.NewReplacer(mSlice...)
			x := r.Replace(source)
			if x != resultStd {
				b.Fatalf("Unexpected result\n%q\nExpected\n%q\n", x, resultStd)
			}
		}
	})
}

func BenchmarkTextTemplate(b *testing.B) {
	s := strings.Replace(source, "{{", "{{.", -1)
	t, err := template.New("test").Parse(s)
	if err != nil {
		b.Fatalf("Error when parsing template: %s", err)
	}

	mm := make(map[string]string)
	for k, v := range m {
		mm[k] = string(v.([]byte))
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var w bytes.Buffer
		for pb.Next() {
			if err := t.Execute(&w, mm); err != nil {
				b.Fatalf("error when executing template: %s", err)
			}
			x := w.Bytes()
			if !bytes.Equal(x, resultTextTemplateBytes) {
				b.Fatalf("unexpected result\n%q\nExpected\n%q\n", x, resultTextTemplateBytes)
			}
			w.Reset()
		}
	})
}

func BenchmarkExecuteFunc(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var bb bytes.Buffer
		for pb.Next() {
			ExecuteFunc(source, "{{", "}}", &bb, testTagFunc)
			bb.Reset()
		}
	})
}

func testTagFunc(w io.Writer, tag string) (int, error) {
	if t, ok := m[tag]; ok {
		return w.Write(t.([]byte))
	}
	return 0, nil
}
