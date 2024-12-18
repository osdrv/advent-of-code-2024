package main

import (
	"fmt"
	"sort"
)

func readPoint(s string) Point2 {
	var p Point2
	_, err := fmt.Sscanf(s, "%d,%d", &p.x, &p.y)
	noerr(err)
	return p
}

const (
	OBST = 1
)

func shortestPath(S, E Point2, B []Point2) int {
	F := makeIntField(E.y+1, E.x+1)
	Q := make([]Point3, 0, 1)
	V := make(map[Point2]int)
	Q = append(Q, Point3{S.x, S.y, 0})

	for _, b := range B {
		F[b.y][b.x] = OBST
	}

	minpath := ALOT
	var H Point3
	for len(Q) > 0 {
		H, Q = Q[0], Q[1:]
		debugf("H: %+v", H)
		H2 := Point2{H.x, H.y}
		if prev, ok := V[H2]; ok && prev <= H.z {
			continue
		}
		V[H2] = H.z
		if H.x == E.x && H.y == E.y {
			debugf("reached end: %d", H.z)
			minpath = min(minpath, H.z)
			continue
		}
		for _, step := range STEPS4 {
			nx, ny := H.x+step[0], H.y+step[1]
			if nx < 0 || nx >= len(F[0]) || ny < 0 || ny >= len(F) {
				continue
			}
			if F[ny][nx] == OBST {
				continue
			}
			NP2 := Point2{nx, ny}
			NP3 := Point3{nx, ny, H.z + 1}
			if prev, ok := V[NP2]; !ok || prev > NP3.z {
				Q = append(Q, NP3)
			}
		}
	}

	return minpath
}

func main() {
	lines := input()

	S := Point2{0, 0}
	ix := 0
	E := readPoint(lines[ix])
	ix++
	nBytes := parseInt(lines[ix])
	ix += 2

	B := make([]Point2, 0, len(lines)-3)
	for ix < len(lines) {
		B = append(B, readPoint(lines[ix]))
		ix++
	}
	debugf("B: %+v", B)

	res1 := shortestPath(S, E, B[:nBytes])
	printf("res1: %d", res1)

	ix = sort.Search(len(B)-nBytes, func(i int) bool {
		return shortestPath(S, E, B[:nBytes+i]) == ALOT
	})
	p := B[nBytes+ix-1]
	printf("res2: %d,%d", p.x, p.y)
}
