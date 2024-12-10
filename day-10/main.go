package main

import "fmt"

func countTrails(F [][]byte, head Point2) (int, int) {
	var rec func(Point2) (int, int)
	U := make(map[Point2]bool)
	rec = func(p Point2) (int, int) {
		score, rating := 0, 0
		v := F[p.y][p.x]
		if v == 9 {
			if _, ok := U[p]; !ok {
				U[p] = true
				return 1, 1
			}
			return 0, 1
		}
		for _, step := range STEPS4 {
			nx := p.x + step[0]
			ny := p.y + step[1]
			if nx < 0 || nx >= len(F[0]) || ny < 0 || ny >= len(F) {
				continue
			}
			nv := F[ny][nx]
			if nv-v != 1 {
				continue
			}
			nscore, nrating := rec(Point2{x: nx, y: ny})
			score += nscore
			rating += nrating
		}
		return score, rating
	}

	return rec(head)
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	heads := make([]Point2, 0, 1)
	for i, line := range lines {
		row := make([]byte, 0, len(line))
		for j, c := range line {
			v := byte(c) - '0'
			row = append(row, v)
			if v == 0 {
				heads = append(heads, Point2{x: j, y: i})
			}
		}
		F = append(F, row)
	}

	debugf("heads: %+v", heads)

	fmt.Println(printNumField(F, ""))

	sum1 := 0
	sum2 := 0
	for _, head := range heads {
		score, rating := countTrails(F, head)
		printf("trail head at %+v has a score %d", head, score)
		sum1 += score
		sum2 += rating
	}

	printf("sum1: %d", sum1)
	printf("sum2: %d", sum2)
}
