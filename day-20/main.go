package main

const (
	START = 'S'
	END   = 'E'
	WALL  = '#'
	EMPTY = '.'
)

type QItem struct {
	P    Point2
	D    int
	PATH []Point2
}

func traverse(F [][]byte, S, E Point2) (int, []Point2) {
	V := make(map[Point2]int)
	Q := make([]QItem, 0, 1)
	Q = append(Q, QItem{
		P:    S,
		D:    0,
		PATH: []Point2{S},
	})
	var H QItem

	mindist := ALOT
	var minpath []Point2
	for len(Q) > 0 {
		H, Q = Q[0], Q[1:]
		P2 := Point2{x: H.P.x, y: H.P.y}
		if prev, ok := V[P2]; ok && prev <= H.D {
			continue
		}

		V[P2] = H.D
		if H.P == E {
			debugf("reached end in %d", H.D)
			if mindist > H.D {
				mindist = H.D
				minpath = H.PATH
			}
			continue
		}

		for _, s := range STEPS4 {
			nx, ny := H.P.x+s[0], H.P.y+s[1]
			if nx < 0 || nx >= len(F[0]) || ny < 0 || ny >= len(F) {
				continue
			}
			if F[ny][nx] == WALL {
				continue
			}
			P2 := Point2{x: nx, y: ny}
			D := H.D + 1
			if prev, ok := V[P2]; !ok || prev > D {
				path := make([]Point2, len(H.PATH)+1)
				copy(path, H.PATH)
				path[len(path)-1] = Point2{x: nx, y: ny}
				qi := QItem{
					P:    Point2{x: nx, y: ny},
					D:    D,
					PATH: path,
				}
				Q = append(Q, qi)
			}
		}
	}

	return mindist, minpath
}

func computeCands(DIST map[Point2]int, p Point2, R int) []Point2 {
	cands := make([]Point2, 0, 1)
	for i := p.y - R; i <= p.y+R; i++ {
		for j := p.x - R; j <= p.x+R; j++ {
			dist := abs(p.y-i) + abs(p.x-j)
			if dist > 0 && dist <= R {
				cand := Point2{x: j, y: i}
				if d, ok := DIST[cand]; ok && d > DIST[p] {
					cands = append(cands, cand)
				}
			}
		}
	}
	return cands
}

func countShortcuts(P []Point2, R int) map[int]int {
	DIST := make(map[Point2]int)
	for dist, p := range P {
		DIST[p] = dist
	}

	DISTCNT := make(map[int]int)
	for dist, p := range P {
		cands := computeCands(DIST, p, R)
		for _, cand := range cands {
			straight := abs(p.x-cand.x) + abs(p.y-cand.y)
			dp := DIST[cand] - dist
			if dp > 0 && straight < dp {
				DISTCNT[dp-straight]++
			}
		}
	}

	return DISTCNT
}

func countDistAbove(DIST map[int]int, R int) int {
	res := 0
	for dist, cnt := range DIST {
		if dist >= R {
			res += cnt
		}
	}
	return res
}

func main() {
	lines := input()
	F := make([][]byte, 0, 1)
	var S, E Point2
	for i, line := range lines {
		F = append(F, []byte(line))
		for j := 0; j < len(lines[i]); j++ {
			if F[i][j] == START {
				S = Point2{x: j, y: i}
			} else if F[i][j] == END {
				E = Point2{x: j, y: i}
			}
		}
	}

	d, path := traverse(F, S, E)
	debugf("nocheats: %d", d)
	debugf("path: %+v", path)

	cnt1 := countShortcuts(path, 2)
	debugf("cnt1: %+v", cnt1)
	printf("res1: %d", countDistAbove(cnt1, 100))
	cnt2 := countShortcuts(path, 20)
	debugf("cnt2: %+v", cnt2)
	printf("res2: %d", countDistAbove(cnt2, 100))
}
