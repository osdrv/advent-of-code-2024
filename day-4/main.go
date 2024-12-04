package main

var (
	XMAS = []byte("XMAS")
	MAS  = []byte("MAS")
)

func traverseCnt(F [][]byte, i, j, di, dj int, gr []byte, ix int) int {
	if i < 0 || i >= len(F) || j < 0 || j >= len(F[i]) {
		return 0
	}
	if F[i][j] != gr[ix] {
		return 0
	}
	//debugf("%c at (%d,%d)", gr[ix], i, j)
	if ix == len(gr)-1 {
		return 1
	}

	ni, nj := i+di, j+dj
	if ni < 0 || ni >= len(F) || nj < 0 || nj >= len(F[ni]) {
		return 0
	}
	if F[ni][nj] == gr[ix+1] {
		return traverseCnt(F, ni, nj, di, dj, gr, ix+1)
	}
	return 0
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	for _, line := range lines {
		F = append(F, []byte(line))
	}

	cnt1 := 0
	for i := 0; i < len(F); i++ {
		for j := 0; j < len(F[i]); j++ {
			if F[i][j] == XMAS[0] {
				for _, s := range STEPS8 {
					cnt1 += traverseCnt(F, i, j, s[0], s[1], XMAS, 0)
				}
			}
		}
	}

	printf("cnt: %d", cnt1)

	cnt2 := 0

	for i := 0; i < len(F); i++ {
		for j := 0; j < len(F[i]); j++ {
			if F[i][j] == MAS[0] {
				for _, s := range [][2]int{{1, 1}, {-1, -1}, {1, -1}, {-1, 1}} {
					if traverseCnt(F, i, j, s[0], s[1], MAS, 0) > 0 {
						debugf("candidate at (%d,%d) step (%d,%d)", i, j, s[0], s[1])
						if traverseCnt(F, i, j+2*s[1], s[0], -1*s[1], MAS, 0) > 0 {
							debugf("match at (%d,%d) step (%d,%d)", i, j+2*s[1], s[0], -1*s[1])
							cnt2++
						}
						if traverseCnt(F, i+2*s[0], j, -1*s[0], s[1], MAS, 0) > 0 {
							debugf("match at (%d,%d) step (%d,%d)", i+2*s[0], j, -1*s[0], s[1])
							cnt2++
						}
					}
				}
			}
		}
	}
	// TODO: I'm overcounting it by 2
	cnt2 /= 2

	printf("cnt2: %d", cnt2)
}
