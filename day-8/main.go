package main

func computeAntinodes(F [][]byte, aa []Point2, steps int) []Point2 {
	AN := make(map[Point2]bool)
	for i := 0; i < len(aa); i++ {
		for j := i + 1; j < len(aa); j++ {
			dx := aa[j].x - aa[i].x
			dy := aa[j].y - aa[i].y
			debugf("dx: %d dy: %d", dx, dy)

			k := 0
			p1, p2 := aa[i], aa[j]
			if steps < 0 {
				AN[p1] = true
				AN[p2] = true
			}
			for (steps > 0 && k < steps) || (steps < 0) {
				p1 = Point2{x: p1.x - dx, y: p1.y - dy}
				if p1.x < 0 || p1.x >= len(F[0]) || p1.y < 0 || p1.y >= len(F) {
					break
				}
				AN[p1] = true
				k++
			}
			k = 0
			for (steps > 0 && k < steps) || (steps < 0) {
				p2 = Point2{x: p2.x + dx, y: p2.y + dy}
				if p2.x < 0 || p2.x >= len(F[0]) || p2.y < 0 || p2.y >= len(F) {
					break
				}
				AN[p2] = true
				k++
			}
		}
	}
	res := make([]Point2, 0, len(AN))
	for p := range AN {
		res = append(res, p)
	}
	return res
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	A := make(map[byte][]Point2)
	for i, line := range lines {
		F = append(F, []byte(line))
		for j := 0; j < len(F[i]); j++ {
			if c := F[i][j]; c != '.' {
				A[c] = append(A[c], Point2{x: j, y: i})
			}
		}
	}

	AN := make(map[Point2]bool)
	AN2 := make(map[Point2]bool)

	for c, aa := range A {
		if len(aa) < 2 {
			debugf("There is only 1 antenna %c, no antinodes", c)
			continue
		}
		an := computeAntinodes(F, aa, 1)
		debugf("Antinodes for %c: %+v", c, an)
		for _, p := range an {
			AN[p] = true
		}
		an2 := computeAntinodes(F, aa, -1)
		debugf("Antinodes for %c: %+v", c, an2)
		for _, p := range an2 {
			AN2[p] = true
		}
	}

	printf("Antinodes: %d", len(AN))
	printf("Antinodes2: %d", len(AN2))
}
