package main

import (
	"regexp"
	"strconv"
	"strings"
)

func parseMulExpr(s string) (int, int) {
	ss := strings.SplitN(s[4:len(s)-1], ",", 2)
	a, err := strconv.Atoi(ss[0])
	noerr(err)
	b, err := strconv.Atoi(ss[1])
	noerr(err)

	return a, b
}

const (
	DO   = "do()"
	DONT = "don't()"
)

func main() {
	lines := input()
	debugf("file data: %+v", lines)

	re := regexp.MustCompile(`mul\((\d+)\,(\d+)\)|do\(\)|don't\(\)`)

	sum1 := 0
	sum2 := 0
	do := true
	for _, line := range lines {
		found := re.FindAllString(line, -1)
		for _, f := range found {
			if f == DO {
				do = true
			} else if f == DONT {
				do = false
			} else {
				a, b := parseMulExpr(f)
				sum1 += a * b
				if do {
					sum2 += a * b
				}
			}
		}
	}
	printf("sum1: %d", sum1)
	printf("sum2: %d", sum2)
}
