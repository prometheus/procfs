package util

import (
	"strconv"
	"strings"
)

// TODO(mdlayher): util packages are an anti-pattern and this should be moved
// somewhere else that is more focused in the future.

// A ValueParser enables parsing a single string into a variety of data types
// in a concise and safe way. The Err method must be invoked after invoking
// any other methods to ensure a value was successfully parsed.
type ValueParser struct {
	v   string
	err error
}

// NewValueParser creates a ValueParser using the input string.
func NewValueParser(v string) *ValueParser {
	return &ValueParser{v: v}
}

// PInt64 interprets the underlying value as an int64 and returns a pointer to
// that value.
func (vp *ValueParser) PInt64() *int64 {
	if vp.err != nil {
		return nil
	}

	var (
		base = 10
		in   = vp.v
	)

	// Is this value stored in hexadecimal instead of decimal?
	if strings.HasPrefix(vp.v, "0x") {
		base = 16
		in = strings.TrimPrefix(in, "0x")
	}

	v, err := strconv.ParseInt(in, base, 64)
	if err != nil {
		vp.err = err
		return nil
	}

	return &v
}

// Err returns the last error, if any, encountered by the ValueParser.
func (vp *ValueParser) Err() error {
	return vp.err
}
