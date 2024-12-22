package main

func computeSecret(n int, rounds int) (int, [][2]int) {
	prev := n % 10
	delta := make([][2]int, 0, 1)
	for i := 0; i < rounds; i++ {
		n = (n ^ (n * 64)) % 16777216
		n = (n ^ (n / 32)) % 16777216
		n = (n ^ (n * 2048)) % 16777216
		nmod10 := n % 10
		delta = append(delta, [2]int{nmod10, nmod10 - prev})
		prev = nmod10
	}
	return n, delta
}

func getSequence(delta [][2]int) map[[4]int]int {
	rmap := make(map[[4]int]int)
	for i := 0; i < len(delta)-4; i++ {
		k := [4]int{delta[i][1], delta[i+1][1], delta[i+2][1], delta[i+3][1]}
		if _, ok := rmap[k]; ok {
			continue
		}
		rmap[k] = delta[i+3][0]
	}
	return rmap
}

func main() {
	lines := input()

	res1 := 0
	seqkeys := make(map[[4]int]int)
	for _, line := range lines {
		n := parseInt(line)
		s, delta := computeSecret(n, 2000)
		debugf("num=%d, s=%d", n, s)
		debugf("delta: %+v", delta[:10])
		seqs := getSequence(delta)
		for seq, price := range seqs {
			seqkeys[seq] += price
		}
		//debugf("rmap: %+v", rmap)
		debugf("buyer %d: %d", n, seqs[[4]int{-2, 1, -1, 3}])
		res1 += s
	}
	printf("res1: %d", res1)

	maxsum := -ALOT
	var maxseq [4]int
	for seq, sum := range seqkeys {
		if sum > maxsum {
			maxsum = sum
			maxseq = seq
		}
	}
	debugf("should be: %d", seqkeys[[4]int{-2, 1, -1, 3}])
	debugf("maxseq: %+v", maxseq)
	printf("res2=%d", maxsum)
}
