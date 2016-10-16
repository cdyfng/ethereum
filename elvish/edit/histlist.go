package edit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/elves/elvish/store"
	"github.com/elves/elvish/util"
)

// Command history listing mode.

var ErrStoreOffline = errors.New("store offline")

type histlist struct {
	listing
	all      []string
	filtered []string
}

func (hl *histlist) Len() int {
	return len(hl.filtered)
}

func (hl *histlist) Show(i, width int) styled {
	entry := hl.filtered[i]
	return unstyled(util.TrimEachLineWcwidth(entry, width))
}

func (hl *histlist) Filter(filter string) int {
	hl.filtered = nil
	for _, item := range hl.all {
		if strings.Contains(item, filter) {
			hl.filtered = append(hl.filtered, item)
		}
	}
	// Select the last entry.
	return len(hl.filtered) - 1
}

func (hl *histlist) Accept(i int, ed *Editor) {
	line := hl.filtered[i]
	if len(ed.line) > 0 {
		line = "\n" + line
	}
	ed.insertAtDot(line)
}

func (hl *histlist) ModeTitle(i int) string {
	return fmt.Sprintf(" HISTORY #%d ", i)
}

func startHistlist(ed *Editor) {
	hl, err := newHistlist(ed.store)
	if err != nil {
		ed.Notify("%v", err)
		return
	}

	ed.histlist = hl
	// ed.histlist = newListing(modeHistoryListing, hl)
	ed.mode = ed.histlist
}

func newHistlist(s *store.Store) (*histlist, error) {
	if s == nil {
		return nil, ErrStoreOffline
	}
	seq, err := s.NextCmdSeq()
	if err != nil {
		return nil, err
	}
	all, err := s.Cmds(0, seq)
	if err != nil {
		return nil, err
	}
	hl := &histlist{all: all}
	hl.listing = newListing(modeHistoryListing, hl)
	return hl, nil
}
