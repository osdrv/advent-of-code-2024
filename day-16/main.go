package main

const (
	START = 'S'
	END   = 'E'
	WALL  = '#'
	EMPTY = '.'
)

const (
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

var (
	STEPS = map[int]Point2{
		NORTH: {0, -1},
		EAST:  {1, 0},
		SOUTH: {0, 1},
		WEST:  {-1, 0},
	}
)

const (
	SCORE_MOV = 1
	SCORE_ROT = 1000
)

type Point4 struct {
	x, y, z int
	s       int
	h       []int
}

func (p Point4) Point3() Point3 {
	return Point3{p.x, p.y, p.z}
}

func shortestPathScore(F [][]byte, start, end Point2) (int, map[int]bool) {
	q := make([]Point4, 0, 1)
	sp := Point4{
		x: start.x,
		y: start.y,
		z: EAST,
		s: 0,
		h: []int{start.y<<32 | start.x},
	}
	q = append(q, sp)
	V := make(map[Point3]int)
	V[sp.Point3()] = sp.s

	BT := make(map[int]bool)

	minScore := ALOT
	var head Point4
	for len(q) > 0 {
		head, q = q[0], q[1:]
		debugf("head: %+v", head)
		if prev, ok := V[head.Point3()]; ok && prev < head.s {
			debugf("%+v", V)
			debugf("skip because we were here before with a lower score")
			continue
		}
		score := head.s
		if end.x == head.x && end.y == head.y {
			debugf("reached end: %d", score)
			if minScore < score {
				continue
			}
			if minScore > score {
				BT = make(map[int]bool)
			}
			minScore = min(minScore, score)
			for _, h := range head.h {
				BT[h] = true
			}
			continue
		}

		cand1 := head
		cand1.z += 1
		cand1.z %= 4
		cand1.s += SCORE_ROT
		if prev, ok := V[cand1.Point3()]; !ok || prev >= cand1.s {
			V[cand1.Point3()] = cand1.s
			q = append(q, cand1)
		}

		cand2 := head
		cand2.z -= 1
		cand2.z += 4
		cand2.z %= 4
		cand2.s += SCORE_ROT
		if prev, ok := V[cand2.Point3()]; !ok || prev >= cand2.s {
			V[cand2.Point3()] = cand2.s
			q = append(q, cand2)
		}

		cand3 := head
		cand3.z += 2
		cand3.z %= 4
		cand3.s += 2 * SCORE_ROT
		if prev, ok := V[cand3.Point3()]; !ok || prev >= cand3.s {
			V[cand3.Point3()] = cand3.s
			q = append(q, cand3)
		}

		cand4 := head
		cand4.x += STEPS[cand4.z].x
		cand4.y += STEPS[cand4.z].y
		cand4.s += SCORE_MOV
		if F[cand4.y][cand4.x] != WALL {
			if prev, ok := V[cand4.Point3()]; !ok || prev >= cand4.s {
				V[cand4.Point3()] = cand4.s
				nh := make([]int, len(head.h))
				copy(nh, head.h)
				nh = append(nh, cand4.y<<32|cand4.x)
				cand4.h = nh
				q = append(q, cand4)
			}
		}
	}

	return minScore, BT
}

func printTiles(tiles map[int]bool) {
	for p := range tiles {
		x, y := p&0xFFFFFFFF, p>>32
		debugf("t: {x: %d, y: %d}", x, y)
	}
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	var start, end Point2
	for i, line := range lines {
		F = append(F, []byte(line))
		for j := 0; j < len(lines[i]); j++ {
			if lines[i][j] == START {
				start = Point2{x: j, y: i}
			} else if lines[i][j] == END {
				end = Point2{x: j, y: i}
			}
		}
	}

	debugf("start: %+v, end: %+v", start, end)

	res1, tiles := shortestPathScore(F, start, end)
	printf("res1: %d", res1)
	res2 := len(tiles)

	printTiles(tiles)

	printf("res2: %d", res2)
}
