package main

func halves(n uint64) (uint64, uint64) {
	nd := numDigs(n)
	a, b := uint64(n), uint64(0)
	bf := uint64(1)
	for i := 0; i < nd/2; i++ {
		b += (a % 10) * bf
		bf *= 10
		a /= 10
	}
	return a, b
}

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

const (
	//NBLINKS = 75
	NBLINKS = 25
)

func main() {
	lines := input()
	stones := parseInts(lines[0])

	sum1 := evolve(stones, 25)
	printf("sum1: %d", sum1)

	sum2 := evolve(stones, 75)
	printf("sum2: %d", sum2)
}
