package main

func sign(n int) int {
	if n < 0 {
		return -1
	} else if n > 0 {
		return 1
	}
	return 0
}

func isSafe(ns []int, least, most int) bool {
	s := sign(ns[1] - ns[0])
	if s == 0 {
		// strictly increasing or decreasing
		return false
	}
	for i := 1; i < len(ns); i++ {
		delta := abs(ns[i] - ns[i-1])
		if delta < least || delta > most {
			return false
		}
		ns := sign(ns[i] - ns[i-1])
		if ns != s {
			return false
		}
	}
	return true
}

func main() {
	lines := input()
	debugf("file data: %+v", lines)
	nums := make([][]int, 0, len(lines))
	for _, line := range lines {
		nums = append(nums, parseInts(line))
	}

	debugf("nums: %+v", nums)

	safe := 0
	for _, ns := range nums {
		if isSafe(ns, 1, 3) {
			safe++
		}
	}

	debugf("safe: %d", safe)

	safe2 := 0
NEXT:
	for _, ns := range nums {
		if isSafe(ns, 1, 3) {
			safe2++
			continue
		}

		for i := 0; i < len(ns); i++ {
			nns := make([]int, 0, len(ns)-1)
			nns = append(nns, ns[:i]...)
			nns = append(nns, ns[i+1:]...)
			if isSafe(nns, 1, 3) {
				debugf("level %+v is safe by removing level %d", ns, ns[i])
				safe2++
				continue NEXT
			}
		}
	}
	debugf("safe2: %d", safe2)
}
