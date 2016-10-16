package edit

import (
	"sort"

	"github.com/elves/elvish/eval"
	"github.com/elves/elvish/parse"
)

// styled is a piece of text with style.
type styled struct {
	text  string
	style string
}

func unstyled(s string) styled {
	return styled{s, ""}
}

func (s *styled) addStyle(st string) {
	s.style = joinStyle(s.style, st)
}

func (s *styled) Kind() string {
	return "styled"
}

func (s *styled) String() string {
	return "\033[" + s.style + "m" + s.text + "\033[m"
}

func (s *styled) Repr(indent int) string {
	return "(le:styled " + parse.Quote(s.text) + " " + parse.Quote(s.style) + ")"
}

func styledBuiltin(ec *eval.EvalCtx, text, style string) {
	out := ec.OutputChan()
	out <- &styled{text, style}
}

// Boilerplates for sorting.

type styleds []styled

func (s styleds) Len() int           { return len(s) }
func (s styleds) Less(i, j int) bool { return s[i].text < s[j].text }
func (s styleds) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func sortStyleds(s []styled) {
	sort.Sort(styleds(s))
}
