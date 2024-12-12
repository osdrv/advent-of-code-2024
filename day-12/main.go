package main

import "sort"

const (
	HOR = 0
	VER = 1
)

func sides(reg []Point2) int {
	e := edges(reg)
	em := make(map[Point3]bool)
	for _, p := range e {
		em[p] = true
	}

	vis := make(map[Point3]bool)
	var recurse func(p Point3)
	recurse = func(p Point3) {
		vis[p] = true
		if p.z == HOR {
			for _, dx := range []int{1, -1} {
				np := Point3{x: p.x + dx, y: p.y, z: p.z}
				npy1, npy2 := np, np
				npy1.z = VER
				npy2.z = VER
				npy1.y--
				if em[npy1] && em[npy2] {
					// cross
					continue
				}
				if _, ok := em[np]; ok && !vis[np] {
					recurse(np)
				}
			}
		} else {
			for _, dy := range []int{1, -1} {
				np := Point3{x: p.x, y: p.y + dy, z: p.z}
				npy1, npy2 := np, np
				npy1.z = HOR
				npy2.z = HOR
				npy1.x--
				if em[npy1] && em[npy2] {
					// cross
					continue
				}
				if _, ok := em[np]; ok && !vis[np] {
					recurse(np)
				}
				if _, ok := em[np]; ok && !vis[np] {
					recurse(np)
				}
			}
		}
	}

	sort.Slice(e, func(a, b int) bool {
		return (e[a].x < e[b].x) || (e[a].x == e[b].x && e[a].y < e[b].y)
	})

	s := 0
	for _, e := range e {
		if vis[e] {
			continue
		}
		recurse(e)
		s++
	}

	return s
}

func edges(reg []Point2) []Point3 {
	E := make(map[Point3]int)

	for _, p := range reg {
		E[Point3{x: p.x, y: p.y, z: HOR}]++     // upper edge
		E[Point3{x: p.x, y: p.y, z: VER}]++     // left edge
		E[Point3{x: p.x, y: p.y + 1, z: HOR}]++ // bottom edge
		E[Point3{x: p.x + 1, y: p.y, z: VER}]++ // right edge
	}

	res := make([]Point3, 0, 1)
	for e, v := range E {
		if v == 1 {
			res = append(res, e)
		}
	}
	return res
}

func perim(reg []Point2) int {
	e := edges(reg)
	return len(e)
}

func square(reg []Point2) int {
	return len(reg)
}

func findRegions(F [][]byte) [][]Point2 {
	v := make([][]bool, len(F))
	for i := 0; i < len(F); i++ {
		v[i] = make([]bool, len(F[i]))
	}

	regions := make([][]Point2, 0)

	var flood func(i, j int) []Point2
	flood = func(i, j int) []Point2 {
		v[i][j] = true
		res := make([]Point2, 0, 1)
		res = append(res, Point2{y: i, x: j})
		for _, s := range STEPS4 {
			ni, nj := i+s[0], j+s[1]
			if ni < 0 || ni >= len(F) || nj < 0 || nj >= len(F[ni]) {
				continue
			}
			if v[ni][nj] || F[i][j] != F[ni][nj] {
				continue
			}
			res = append(res, flood(ni, nj)...)
		}
		return res
	}

	for i := 0; i < len(F); i++ {
		for j := 0; j < len(F[i]); j++ {
			if v[i][j] {
				continue
			}
			regions = append(regions, flood(i, j))
		}
	}

	return regions
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	for _, line := range lines {
		F = append(F, []byte(line))
	}

	regs := findRegions(F)

	debugf("regs: %v", regs)

	sum1 := 0
	sum2 := 0
	for _, reg := range regs {
		p, s, sd := perim(reg), square(reg), sides(reg)
		debugf("Region: %v, perim: %d, square: %d, sides: %d", reg, p, s, sd)
		sum1 += p * s
		sum2 += sd * s
	}

	printf("sum1: %d", sum1)
	printf("sum2: %d", sum2)
}
