package main

import "sort"

func main() {
	lines := input()
	debugf("file data: %+v", lines)

	left, right := make([]int, 0, len(lines)), make([]int, 0, len(lines))
	cntR := make(map[int]int)
	for _, line := range lines {
		nn := parseInts(line)
		left = append(left, nn[0])
		right = append(right, nn[1])
		cntR[nn[1]]++
	}

	debugf("Nums: %+v, %+v", left, right)
	sort.Ints(left)
	sort.Ints(right)

	dist := 0
	sim := 0
	for i := 0; i < len(left); i++ {
		dist += abs(left[i] - right[i])
		sim += left[i] * cntR[left[i]]
	}

	debugf("dist: %d", dist)
	debugf("sim: %d", sim)
}
