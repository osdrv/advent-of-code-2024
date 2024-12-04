package main

import (
	"strconv"
)

type Matcher interface {
	Match(s string, off int) (any, int, bool)
}

type StringMatcher struct {
	s string
}

var _ Matcher = (*StringMatcher)(nil)

func (m *StringMatcher) Match(s string, off int) (any, int, bool) {
	ix := 0
	for ix+off < len(s) && ix < len(m.s) {
		if s[ix+off] != m.s[ix] {
			return "", off + 1, false
		}
		ix++
	}
	if ix == len(m.s) {
		return s[off : off+ix], off + ix, true
	}
	return "", off + 1, false
}

type NumberMatcher struct{}

var _ Matcher = (*NumberMatcher)(nil)

func (m *NumberMatcher) Match(s string, off int) (any, int, bool) {
	ix := 0
	for ix+off < len(s) {
		if !isNumber(s[ix+off]) {
			break
		}
		ix++
	}
	if ix > 0 {
		if n, err := strconv.Atoi(s[off : off+ix]); err == nil {
			return n, off + ix, true
		}
	}
	return 0, off + ix, false
}

type AndMatcher struct {
	mchrz []Matcher
}

var _ Matcher = (*AndMatcher)(nil)

func (m *AndMatcher) Match(s string, off int) (any, int, bool) {
	res := make([]any, 0, len(m.mchrz))
	for _, mchr := range m.mchrz {
		any, newOff, ok := mchr.Match(s, off)
		if !ok {
			return nil, newOff, false
		}
		res = append(res, any)
		off = newOff
	}
	return res, off, true
}

type OrMatcher struct {
	mchrz []Matcher
}

var _ Matcher = (*OrMatcher)(nil)

func (m *OrMatcher) Match(s string, off int) (any, int, bool) {
	for _, mchr := range m.mchrz {
		any, newOff, ok := mchr.Match(s, off)
		if ok {
			return any, newOff, true
		}
	}
	return nil, off + 1, false
}

const (
	DO   = "do()"
	DONT = "don't()"
)

var (
	MTCH_DO   = &StringMatcher{DO}
	MTCH_DONT = &StringMatcher{DONT}
	MTCH_MUL  = &AndMatcher{
		mchrz: []Matcher{
			&StringMatcher{"mul("},
			&NumberMatcher{},
			&StringMatcher{","},
			&NumberMatcher{},
			&StringMatcher{")"},
		},
	}
	MTCH_ALL = &OrMatcher{
		mchrz: []Matcher{
			MTCH_DO,
			MTCH_DONT,
			MTCH_MUL,
		},
	}
)

// I don't always regex, but when I do, I use my own regex engine
func main() {
	lines := input()

	sum1 := 0
	sum2 := 0
	do := true
	for _, line := range lines {
		ix := 0
		for ix < len(line) {
			found, newIx, ok := MTCH_ALL.Match(line, ix)
			if ok {
				debugf("Found: %+v", found)
				if found == DO {
					do = true
				} else if found == DONT {
					do = false
				} else {
					ff := found.([]any)
					a, b := ff[1].(int), ff[3].(int)
					sum1 += a * b
					if do {
						sum2 += a * b
					}
				}
			}
			ix = newIx
		}
	}
	printf("sum1: %d", sum1)
	printf("sum2: %d", sum2)
}
