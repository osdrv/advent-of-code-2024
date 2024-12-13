package main

import (
	"fmt"
	"math"
)

type Point2_64 struct {
	x, y int64
}

type Game struct {
	A, B  Point2_64
	Prize Point2_64
}

const (
	COST_A = int64(3)
	COST_B = int64(1)
)

func play(g Game, maxpress int64) (int64, bool) {
	/*
		px = ax*j + bx*k
		py = ay*j + by*k

		k = (px - ax * j) / bx
		py = ay * j + (px - ax * j) * by / bx
		py * bx / by = j * ay * bx / by + px - ax * j
		py * bx / by = j * (ay * bx / by - ax) + px
		j = (py * bx / by - px) / (ay * bx / by - ax)

	*/
	ax, ay := float64(g.A.x), float64(g.A.y)
	bx, by := float64(g.B.x), float64(g.B.y)
	px, py := float64(g.Prize.x), float64(g.Prize.y)

	j := (py*bx/by - px) / (ay*bx/by - ax)
	k := (px - ax*j) / bx

	if j < 0 || k < 0 {
		return -1, false
	}

	ji, ki := int64(math.Round(j)), int64(math.Round(k))
	debugf("j: %f(%d), k:%f(%d)", j, ji, k, ki)

	iswin := ji*g.A.x+ki*g.B.x == g.Prize.x &&
		ji*g.A.y+ki*g.B.y == g.Prize.y &&
		ji <= maxpress &&
		ki <= maxpress

	return ji*COST_A + ki*COST_B, iswin
}

func main() {
	lines := input()
	ix := 0
	games := make([]Game, 0, 1)
	for ix < len(lines) {
		var g Game
		var err error
		_, err = fmt.Sscanf(lines[ix], "Button A: X+%d, Y+%d", &g.A.x, &g.A.y)
		noerr(err)
		_, err = fmt.Sscanf(lines[ix+1], "Button B: X+%d, Y+%d", &g.B.x, &g.B.y)
		noerr(err)
		_, err = fmt.Sscanf(lines[ix+2], "Prize: X=%d, Y=%d", &g.Prize.x, &g.Prize.y)
		noerr(err)
		ix += 4
		games = append(games, g)
	}

	debugf("games: %+v", games)

	sum1 := int64(0)
	for _, g := range games {
		if tok, ok := play(g, 100); ok {
			debugf("game %+v is winnable with %d tokens", g, tok)
			sum1 += tok
		} else {
			debugf("game %+v is not winnable", g)
		}
	}

	sum2 := int64(0)
	OFF := int64(10000000000000)
	for _, g := range games {
		ng := g
		ng.Prize.x += OFF
		ng.Prize.y += OFF

		if tok, ok := play(ng, 99999999999999999); ok {
			debugf("game %+v is winnable with %d tokens", g, tok)
			sum2 += tok
		} else {
			debugf("game %+v is not winnable", g)
		}
	}

	printf("sum2: %d", sum2)
}
