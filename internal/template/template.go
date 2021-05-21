// Package fasttemplate implements simple and fast template library.
//
// Fasttemplate is faster than text/template, strings.Replace
// and strings.Replacer.
//
// Fasttemplate ideally fits for fast and simple placeholders' substitutions.
package template

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/lucasepe/tbd/internal/bytebufferpool"
)

// ExecuteFunc calls f on each template tag (placeholder) occurrence.
//
// Returns the number of bytes written to w.
//
// This function is optimized for constantly changing templates.
// Use Template.ExecuteFunc for frozen templates.
func ExecuteFunc(template, startTag, endTag string, w io.Writer, f TagFunc) (int64, error) {
	s := unsafeString2Bytes(template)
	a := unsafeString2Bytes(startTag)
	b := unsafeString2Bytes(endTag)

	var nn int64
	var ni int
	var err error
	for {
		n := bytes.Index(s, a)
		if n < 0 {
			break
		}
		ni, err = w.Write(s[:n])
		nn += int64(ni)
		if err != nil {
			return nn, err
		}

		s = s[n+len(a):]
		n = bytes.Index(s, b)
		if n < 0 {
			// cannot find end tag - just write it to the output.
			ni, _ = w.Write(a)
			nn += int64(ni)
			break
		}

		tag := strings.TrimSpace(unsafeBytes2String(s[:n]))
		ni, err = f(w, tag)
		nn += int64(ni)
		if err != nil {
			return nn, err
		}
		s = s[n+len(b):]
	}
	ni, err = w.Write(s)
	nn += int64(ni)

	return nn, err
}

// Marks returns the list of all placeholders found in the specified template.
func Marks(template, startTag, endTag string) ([]string, error) {
	list := []string{}
	_, err := ExecuteFunc(template, startTag, endTag, io.Discard,
		func(w io.Writer, tag string) (int, error) {
			return fetchTagFunc(tag, &list)
		})
	return list, err
}

// Execute substitutes template tags (placeholders) with the corresponding
// values from the map m and writes the result to the given writer w.
//
// Substitution map m may contain values with the following types:
//   * []byte - the fastest value type
//   * string - convenient value type
//   * TagFunc - flexible value type
//
// Returns the number of bytes written to w.
//
// This function is optimized for constantly changing templates.
// Use Template.Execute for frozen templates.
func Execute(template, startTag, endTag string, w io.Writer, m map[string]interface{}) (int64, error) {
	return ExecuteFunc(template, startTag, endTag, w,
		func(w io.Writer, tag string) (int, error) {
			return stdTagFunc(w, tag, m)
		})
}

// ExecuteStd works the same way as Execute, but keeps the unknown placeholders.
// This can be used as a drop-in replacement for strings.Replacer
//
// Substitution map m may contain values with the following types:
//   * []byte - the fastest value type
//   * string - convenient value type
//   * TagFunc - flexible value type
//
// Returns the number of bytes written to w.
//
// This function is optimized for constantly changing templates.
// Use Template.ExecuteStd for frozen templates.
func ExecuteStd(template, startTag, endTag string, w io.Writer, m map[string]interface{}) (int64, error) {
	return ExecuteFunc(template, startTag, endTag, w,
		func(w io.Writer, tag string) (int, error) {
			return keepUnknownTagFunc(w, startTag, endTag, tag, m)
		})
}

// ExecuteFuncString calls f on each template tag (placeholder) occurrence
// and substitutes it with the data written to TagFunc's w.
//
// Returns the resulting string that will be empty on error.
func ExecuteFuncString(template, startTag, endTag string, f TagFunc) (string, error) {
	tagsCount := bytes.Count(unsafeString2Bytes(template), unsafeString2Bytes(startTag))
	if tagsCount == 0 {
		return template, nil
	}

	var byteBufferPool bytebufferpool.Pool

	bb := byteBufferPool.Get()
	if _, err := ExecuteFunc(template, startTag, endTag, bb, f); err != nil {
		bb.Reset()
		byteBufferPool.Put(bb)
		return "", err
	}
	s := string(bb.B)
	bb.Reset()
	byteBufferPool.Put(bb)
	return s, nil
}

// ExecuteString substitutes template tags (placeholders) with the corresponding
// values from the map m and returns the result.
//
// Substitution map m may contain values with the following types:
//   * []byte - the fastest value type
//   * string - convenient value type
//   * TagFunc - flexible value type
//
// This function is optimized for constantly changing templates.
// Use Template.ExecuteString for frozen templates.
func ExecuteString(template, startTag, endTag string, m map[string]interface{}) (string, error) {
	return ExecuteFuncString(template, startTag, endTag,
		func(w io.Writer, tag string) (int, error) {
			return stdTagFunc(w, tag, m)
		})
}

// ExecuteStringStd works the same way as ExecuteString, but keeps the unknown placeholders.
// This can be used as a drop-in replacement for strings.Replacer
//
// Substitution map m may contain values with the following types:
//   * []byte - the fastest value type
//   * string - convenient value type
//   * TagFunc - flexible value type
//
// This function is optimized for constantly changing templates.
// Use Template.ExecuteStringStd for frozen templates.
func ExecuteStringStd(template, startTag, endTag string, m map[string]interface{}) (string, error) {
	return ExecuteFuncString(template, startTag, endTag,
		func(w io.Writer, tag string) (int, error) {
			return keepUnknownTagFunc(w, startTag, endTag, tag, m)
		})
}

// TagFunc can be used as a substitution value in the map passed to Execute*.
// Execute* functions pass tag (placeholder) name in 'tag' argument.
//
// TagFunc must be safe to call from concurrently running goroutines.
//
// TagFunc must write contents to w and return the number of bytes written.
type TagFunc func(w io.Writer, tag string) (int, error)

func stdTagFunc(w io.Writer, tag string, m map[string]interface{}) (int, error) {
	v := m[tag]
	if v == nil {
		return 0, nil
	}
	switch value := v.(type) {
	case []byte:
		return w.Write(value)
	case string:
		return w.Write([]byte(value))
	case TagFunc:
		return value(w, tag)
	default:
		return -1, fmt.Errorf("tag=%q contains unexpected value type=%#v. Expected []byte, string or TagFunc", tag, v)
	}
}

func keepUnknownTagFunc(w io.Writer, startTag, endTag, tag string, m map[string]interface{}) (int, error) {
	v, ok := m[tag]
	if !ok {
		if _, err := w.Write(unsafeString2Bytes(startTag)); err != nil {
			return 0, err
		}
		if _, err := w.Write(unsafeString2Bytes(tag)); err != nil {
			return 0, err
		}
		if _, err := w.Write(unsafeString2Bytes(endTag)); err != nil {
			return 0, err
		}
		return len(startTag) + len(tag) + len(endTag), nil
	}
	if v == nil {
		return 0, nil
	}
	switch value := v.(type) {
	case []byte:
		return w.Write(value)
	case string:
		return w.Write([]byte(value))
	case TagFunc:
		return value(w, tag)
	default:
		return -1, fmt.Errorf("tag=%q contains unexpected value type=%#v. Expected []byte, string or TagFunc", tag, v)
	}
}

// fetchTagFunc accumulates all tags in the specified array.
func fetchTagFunc(tag string, arr *[]string) (int, error) {
	*arr = append(*arr, tag)
	return 0, nil
}
