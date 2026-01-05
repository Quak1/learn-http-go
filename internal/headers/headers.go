package headers

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 2, true, nil
	}

	data = bytes.TrimSpace(data[:idx])
	sepIdx := bytes.IndexRune(data, ':')
	if sepIdx == -1 {
		return 0, false, fmt.Errorf("Error: invalid field line, no separator ':' found")
	}
	if !isValidKey(data[:sepIdx]) {
		return 0, false, fmt.Errorf("Error: invalid field line name")
	}

	key := string(data[:sepIdx])
	value := string(bytes.TrimSpace(data[sepIdx+1:]))
	h.Set(key, value)

	return idx + 2, false, nil
}

func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9')
}

func isValidKey(data []byte) bool {
	marks := map[byte]bool{
		'!':  true,
		'#':  true,
		'$':  true,
		'%':  true,
		'&':  true,
		'\'': true,
		'*':  true,
		'+':  true,
		'-':  true,
		'.':  true,
		'^':  true,
		'_':  true,
		'`':  true,
		'|':  true,
		'~':  true,
	}

	for _, b := range data {
		if unicode.IsSpace(rune(b)) {
			return false
		}
		if !isAlphaNumeric(b) && !marks[b] {
			return false
		}
	}

	return true
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	if v, ok := h[key]; ok {
		value = v + ", " + value
	}
	h[key] = value
}

func (h Headers) Get(key string) string {
	lower := strings.ToLower(key)
	return h[lower]
}
