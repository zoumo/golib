package textencoding

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
)

var (
	all   = map[string]encoding.Encoding{}
	alias = map[string]encoding.Encoding{
		"UTF8":   unicode.UTF8,
		"GB2312": simplifiedchinese.HZGB2312,
	}
)

type encodingWithName interface {
	String() string
}

func init() {
	tempall := []encoding.Encoding{}
	tempall = extend(
		tempall,
		simplifiedchinese.All,
		traditionalchinese.All,
		japanese.All,
		korean.All,
		charmap.All,
		unicode.All,
		utf32.All,
	)

	for _, e := range tempall {
		en, ok := e.(encodingWithName)
		if !ok {
			continue
		}
		all[strings.ToUpper(en.String())] = e
	}

	for k, e := range alias {
		all[strings.ToUpper(k)] = e
	}

}

func extend(dest []encoding.Encoding, alls ...[]encoding.Encoding) []encoding.Encoding {
	for _, all := range alls {
		dest = append(dest, all...)
	}
	return dest
}

// IsEncodingSupported checks if the encoding is supported
func IsEncodingSupported(name string) bool {
	_, ok := all[strings.ToUpper(name)]
	return ok
}

// Encode encodes the utf-8 bytes into target encoding
func Encode(s []byte, to string) ([]byte, error) {
	return Transform(s, "UTF-8", to)
}

// Decode decodes the bytes to UTF-8 bytes
func Decode(s []byte, from string) ([]byte, error) {
	return Transform(s, from, "UTF-8")
}

// TransformString decodes the input string with srouce encoding and
// then encodes it into target encoding
func TransformString(s string, from, to string) (string, error) {
	ret, err := Transform([]byte(s), from, to)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

// Transform decodes the input bytes with srouce encoding and
// then encodes them into target encoding
func Transform(s []byte, from, to string) ([]byte, error) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	fromEncoding, ok := all[from]
	if !ok {
		return nil, fmt.Errorf("Unsupported from encoding %v", from)
	}

	toEncoding, ok := all[to]
	if !ok {
		return nil, fmt.Errorf("Unsupported to encoding %v", to)
	}

	reader := transform.NewReader(bytes.NewBuffer(s), transform.Chain(fromEncoding.NewDecoder(), toEncoding.NewEncoder()))

	ret, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
