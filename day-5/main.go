package main

import (
	"fmt"
	"sort"
)

func parseRule(s string) (int, int) {
	var a, b int
	fmt.Sscanf(s, "%d|%d", &a, &b)
	return a, b
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for ix := range a {
		if a[ix] != b[ix] {
			return false
		}
	}
	return true
}

func main() {
	lines := input()
	rules := make(map[int]map[int]bool)

	ix := 0
	for ix < len(lines) {
		if len(lines[ix]) == 0 {
			break
		}
		a, b := parseRule(lines[ix])
		if _, ok := rules[a]; !ok {
			rules[a] = make(map[int]bool)
		}
		rules[a][b] = true
		ix++
	}
	ix++

	debugf("rules: %+v", rules)

	sum1 := 0
	sum2 := 0
	for ix < len(lines) {
		pages := parseInts(lines[ix])
		pagescp := make([]int, len(pages))
		copy(pagescp, pages)
		sort.Slice(pagescp, func(i, j int) bool {
			if rules[pagescp[i]][pagescp[j]] {
				return true
			}
			return false
		})

		if equal(pages, pagescp) {
			debugf("pages %+v is EQUAL", pages)
			sum1 += pages[len(pages)/2]
		} else {
			debugf("pages %+v is NOT EQUAL", pages)
			sum2 += pagescp[len(pagescp)/2]
		}
		ix++
	}

	printf("sum1: %d", sum1)
	printf("sum2: %d", sum2)
}
