package template

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestExecuteFunc(t *testing.T) {
	testExecuteFunc(t, "", "")
	testExecuteFunc(t, "a", "a")
	testExecuteFunc(t, "abc", "abc")
	testExecuteFunc(t, "{foo}", "xxxx")
	testExecuteFunc(t, "a{foo}", "axxxx")
	testExecuteFunc(t, "{foo}a", "xxxxa")
	testExecuteFunc(t, "a{foo}bc", "axxxxbc")
	testExecuteFunc(t, "{foo}{foo}", "xxxxxxxx")
	testExecuteFunc(t, "{foo}bar{foo}", "xxxxbarxxxx")

	// unclosed tag
	testExecuteFunc(t, "{unclosed", "{unclosed")
	testExecuteFunc(t, "{{unclosed", "{{unclosed")
	testExecuteFunc(t, "{un{closed", "{un{closed")

	// test unknown tag
	testExecuteFunc(t, "{unknown}", "zz")
	testExecuteFunc(t, "{foo}q{unexpected}{missing}bar{foo}", "xxxxqzzzzbarxxxx")
}

func testExecuteFunc(t *testing.T, template, expectedOutput string) {
	var bb bytes.Buffer
	ExecuteFunc(template, "{", "}", &bb, func(w io.Writer, tag string) (int, error) {
		if tag == "foo" {
			return w.Write([]byte("xxxx"))
		}
		return w.Write([]byte("zz"))
	})

	output := string(bb.Bytes())
	if output != expectedOutput {
		t.Fatalf("unexpected output for template=%q: %q. Expected %q", template, output, expectedOutput)
	}
}

func TestExecute(t *testing.T) {
	testExecute(t, "", "")
	testExecute(t, "a", "a")
	testExecute(t, "abc", "abc")
	testExecute(t, "{foo}", "xxxx")
	testExecute(t, "a{foo}", "axxxx")
	testExecute(t, "{foo}a", "xxxxa")
	testExecute(t, "a{foo}bc", "axxxxbc")
	testExecute(t, "{foo}{foo}", "xxxxxxxx")
	testExecute(t, "{foo}bar{foo}", "xxxxbarxxxx")

	// unclosed tag
	testExecute(t, "{unclosed", "{unclosed")
	testExecute(t, "{{unclosed", "{{unclosed")
	testExecute(t, "{un{closed", "{un{closed")

	// test unknown tag
	testExecute(t, "{unknown}", "")
	testExecute(t, "{foo}q{unexpected}{missing}bar{foo}", "xxxxqbarxxxx")
}

func testExecute(t *testing.T, template, expectedOutput string) {
	var bb bytes.Buffer
	Execute(template, "{", "}", &bb, map[string]interface{}{"foo": "xxxx"})
	output := string(bb.Bytes())
	if output != expectedOutput {
		t.Fatalf("unexpected output for template=%q: %q. Expected %q", template, output, expectedOutput)
	}
}

func TestExecuteStd(t *testing.T) {
	testExecuteStd(t, "", "")
	testExecuteStd(t, "a", "a")
	testExecuteStd(t, "abc", "abc")
	testExecuteStd(t, "{foo}", "xxxx")
	testExecuteStd(t, "a{foo}", "axxxx")
	testExecuteStd(t, "{foo}a", "xxxxa")
	testExecuteStd(t, "a{foo}bc", "axxxxbc")
	testExecuteStd(t, "{foo}{foo}", "xxxxxxxx")
	testExecuteStd(t, "{foo}bar{foo}", "xxxxbarxxxx")

	// unclosed tag
	testExecuteStd(t, "{unclosed", "{unclosed")
	testExecuteStd(t, "{{unclosed", "{{unclosed")
	testExecuteStd(t, "{un{closed", "{un{closed")

	// test unknown tag
	testExecuteStd(t, "{unknown}", "{unknown}")
	testExecuteStd(t, "{foo}q{unexpected}{missing}bar{foo}", "xxxxq{unexpected}{missing}barxxxx")
}

func testExecuteStd(t *testing.T, template, expectedOutput string) {
	var bb bytes.Buffer
	ExecuteStd(template, "{", "}", &bb, map[string]interface{}{"foo": "xxxx"})
	output := string(bb.Bytes())
	if output != expectedOutput {
		t.Fatalf("unexpected output for template=%q: %q. Expected %q", template, output, expectedOutput)
	}
}

func TestExecuteString(t *testing.T) {
	testExecuteString(t, "", "")
	testExecuteString(t, "a", "a")
	testExecuteString(t, "abc", "abc")
	testExecuteString(t, "{foo}", "xxxx")
	testExecuteString(t, "a{foo}", "axxxx")
	testExecuteString(t, "{foo}a", "xxxxa")
	testExecuteString(t, "a{foo}bc", "axxxxbc")
	testExecuteString(t, "{foo}{foo}", "xxxxxxxx")
	testExecuteString(t, "{foo}bar{foo}", "xxxxbarxxxx")

	// unclosed tag
	testExecuteString(t, "{unclosed", "{unclosed")
	testExecuteString(t, "{{unclosed", "{{unclosed")
	testExecuteString(t, "{un{closed", "{un{closed")

	// test unknown tag
	testExecuteString(t, "{unknown}", "")
	testExecuteString(t, "{foo}q{unexpected}{missing}bar{foo}", "xxxxqbarxxxx")
}

func testExecuteString(t *testing.T, template, expectedOutput string) {
	output, err := ExecuteString(template, "{", "}", map[string]interface{}{"foo": "xxxx"})
	if err != nil {
		t.Fatal(err)
	}

	if output != expectedOutput {
		t.Fatalf("unexpected output for template=%q: %q. Expected %q", template, output, expectedOutput)
	}
}

func TestExecuteStringStd(t *testing.T) {
	testExecuteStringStd(t, "", "")
	testExecuteStringStd(t, "a", "a")
	testExecuteStringStd(t, "abc", "abc")
	testExecuteStringStd(t, "{foo}", "xxxx")
	testExecuteStringStd(t, "a{foo}", "axxxx")
	testExecuteStringStd(t, "{foo}a", "xxxxa")
	testExecuteStringStd(t, "a{foo}bc", "axxxxbc")
	testExecuteStringStd(t, "{foo}{foo}", "xxxxxxxx")
	testExecuteStringStd(t, "{foo}bar{foo}", "xxxxbarxxxx")

	// unclosed tag
	testExecuteStringStd(t, "{unclosed", "{unclosed")
	testExecuteStringStd(t, "{{unclosed", "{{unclosed")
	testExecuteStringStd(t, "{un{closed", "{un{closed")

	// test unknown tag
	testExecuteStringStd(t, "{unknown}", "{unknown}")
	testExecuteStringStd(t, "{foo}q{unexpected}{missing}bar{foo}", "xxxxq{unexpected}{missing}barxxxx")
}

func testExecuteStringStd(t *testing.T, template, expectedOutput string) {
	output, err := ExecuteStringStd(template, "{", "}", map[string]interface{}{"foo": "xxxx"})
	if err != nil {
		t.Fatal(err)
	}

	if output != expectedOutput {
		t.Fatalf("unexpected output for template=%q: %q. Expected %q", template, output, expectedOutput)
	}
}

func expectPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("missing panic")
		}
	}()
	f()
}

func TestExecuteFuncStringWithErr(t *testing.T) {
	var expectErr = errors.New("test111")
	result, err := ExecuteFuncString(`{a} is {b}'s best friend`, "{", "}", func(w io.Writer, tag string) (int, error) {
		if tag == "a" {
			return w.Write([]byte("Alice"))
		}
		return 0, expectErr
	})
	if err != expectErr {
		t.Fatalf("error must be the same as the error returned from f, expect: %s, actual: %s", expectErr, err)
	}
	if result != "" {
		t.Fatalf("result should be an empty string if error occurred")
	}
	result, err = ExecuteFuncString(`{a} is {b}'s best friend`, "{", "}", func(w io.Writer, tag string) (int, error) {
		if tag == "a" {
			return w.Write([]byte("Alice"))
		}
		return w.Write([]byte("Bob"))
	})
	if err != nil {
		t.Fatalf("should success but found err: %s", err)
	}
	if result != "Alice is Bob's best friend" {
		t.Fatalf("expect: %s, but: %s", "Alice is Bob's best friend", result)
	}
}
