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

type QueueItem struct {
	P Point4x
	H []Point2x
}

type Point4 struct {
	x, y, z int
	s       int
	h       []int
}

func (p Point4) Point3() Point3 {
	return Point3{p.x, p.y, p.z}
}

func shortestPathScore(F [][]byte, start, end Point2x) (int, map[Point2x]struct{}) {
	q := make([]QueueItem, 0, 1)
	p := p4(start.x(), start.y(), EAST, 0)
	q = append(q, QueueItem{
		P: p,
		H: []Point2x{start},
	})
	V := make(map[Point3x]int)
	V[p.p3()] = p.s()

	BT := make(map[Point2x]struct{})

	minScore := ALOT
	var head QueueItem
	for len(q) > 0 {
		head, q = q[0], q[1:]
		if prev, ok := V[head.P.p3()]; ok && prev < head.P.s() {
			debugf("skip because we were here before with a lower score")
			continue
		}
		score := head.P.s()
		if end == head.P.p2() {
			debugf("reached end: %d", score)
			if minScore < score {
				continue
			}
			if minScore > score {
				BT = make(map[Point2x]struct{})
			}
			minScore = min(minScore, score)
			for _, h := range head.H {
				BT[h] = struct{}{}
			}
			continue
		}

		cand1 := p4(head.P.x(), head.P.y(), (head.P.z()+1)%4, head.P.s()+SCORE_ROT)
		if prev, ok := V[cand1.p3()]; !ok || prev >= cand1.s() {
			V[cand1.p3()] = cand1.s()
			q = append(q, QueueItem{cand1, head.H})
		}

		cand2 := p4(head.P.x(), head.P.y(), (head.P.z()+3)%4, head.P.s()+SCORE_ROT)
		if prev, ok := V[cand2.p3()]; !ok || prev >= cand2.s() {
			V[cand2.p3()] = cand2.s()
			q = append(q, QueueItem{cand2, head.H})
		}

		cand3 := p4(head.P.x(), head.P.y(), (head.P.z()+2)%4, head.P.s()+2*SCORE_ROT)
		if prev, ok := V[cand3.p3()]; !ok || prev >= cand3.s() {
			V[cand3.p3()] = cand3.s()
			q = append(q, QueueItem{cand3, head.H})
		}

		cand4 := p4(
			head.P.x()+STEPS[head.P.z()].x,
			head.P.y()+STEPS[head.P.z()].y,
			head.P.z(),
			head.P.s()+SCORE_MOV,
		)
		if F[cand4.y()][cand4.x()] != WALL {
			if prev, ok := V[cand4.p3()]; !ok || prev >= cand4.s() {
				V[cand4.p3()] = cand4.s()
				nh := make([]Point2x, len(head.H))
				copy(nh, head.H)
				nh = append(nh, cand4.p2())
				q = append(q, QueueItem{cand4, nh})
			}
		}
	}

	return minScore, BT
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	var start, end Point2x
	for i, line := range lines {
		F = append(F, []byte(line))
		for j := 0; j < len(lines[i]); j++ {
			if lines[i][j] == START {
				start = p2(j, i)
			} else if lines[i][j] == END {
				end = p2(j, i)
			}
		}
	}

	debugf("start: %+v, end: %+v", start, end)

	res1, tiles := shortestPathScore(F, start, end)
	printf("res1: %d", res1)
	res2 := len(tiles)

	printf("res2: %d", res2)
}
