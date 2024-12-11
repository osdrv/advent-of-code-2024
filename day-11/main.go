package main

func evolve(stones []int, n int) uint64 {
	ss := make(map[uint64]uint64)

	for _, s := range stones {
		ss[uint64(s)] += 1
	}

	for step := 0; step < n; step++ {
		ns := make(map[uint64]uint64)
		for num, cnt := range ss {
			if num == 0 {
				ns[1] += cnt
			} else if numDigs(num)%2 == 0 {
				a, b := halves(num)
				ns[a] += cnt
				ns[b] += cnt
			} else {
				ns[num*2024] += cnt
			}
		}
		ss = ns
	}

	sum := uint64(0)
	for _, cnt := range ss {
		sum += cnt
	}
	return sum
}

func main() {
	lines := input()
	stones := parseInts(lines[0])

	sum1 := evolve(stones, 25)
	printf("sum1: %d", sum1)

	sum2 := evolve(stones, 75)
	printf("sum2: %d", sum2)
}
