package eval

import (
	"errors"

	"github.com/elves/elvish/parse"
)

// MapStringString implements MapLike for map[string]string.
type MapStringString map[string]string

var (
	_ MapLike     = MapStringString(nil)
	_ IndexSetter = MapStringString(nil)
)

var (
	ErrValueMustBeString = errors.New("index must be string")
)

func (MapStringString) Kind() string {
	return "map"
}

func (m MapStringString) Repr(indent int) string {
	var builder MapReprBuilder
	builder.Indent = indent
	for k, v := range m {
		builder.WritePair(parse.Quote(k), parse.Quote(v))
	}
	return builder.String()
}

func (m MapStringString) Len() int {
	return len(m)
}

func (m MapStringString) IndexOne(idx Value) Value {
	i, ok := idx.(String)
	if !ok {
		throw(ErrIndexMustBeString)
	}
	v, ok := m[string(i)]
	if !ok {
		throw(errors.New("no such key: " + i.Repr(NoPretty)))
	}
	return String(v)
}

func (m MapStringString) HasKey(idx Value) bool {
	if i, ok := idx.(String); ok {
		if _, ok := m[string(i)]; ok {
			return true
		}
	}
	return false
}

func (m MapStringString) IndexSet(idx Value, val Value) {
	i, ok := idx.(String)
	if !ok {
		throw(ErrIndexMustBeString)
	}
	v, ok := val.(String)
	if !ok {
		throw(ErrValueMustBeString)
	}
	m[string(i)] = string(v)
}
