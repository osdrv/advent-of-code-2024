package main

import "strings"

var MEMO = make(map[string]int)

func traverse2(W string, S []string) int {
	if len(W) == 0 {
		return 1
	}
	if prev, ok := MEMO[W]; ok {
		return prev
	}
	res := 0
	for _, s := range S {
		if strings.HasPrefix(W, s) {
			res += traverse2(W[len(s):], S)
		}
	}
	MEMO[W] = res
	return res
}

func main() {
	lines := input()
	paths := strings.Split(lines[0], ", ")
	des := lines[2:]

	res1 := 0
	res2 := 0
	for _, d := range des {
		debugf("traversing %s", d)
		if n := traverse2(d, paths); n > 0 {
			debugf("design %s is possible in %d ways", d, n)
			res1++
			res2 += n
		} else {
			debugf("design %s is impossible", d)
		}
	}

	printf("res1: %d", res1)
	printf("res2: %d", res2)
}
