package main

import (
	"fmt"
)

type Robot struct {
	P Point2
	V Point2
}

func (r *Robot) String() string {
	return fmt.Sprintf("Robo[P={%d,%d} V={%d,%d}]", r.P.x, r.P.y, r.V.x, r.V.y)
}

func parseRobot(s string) (*Robot, error) {
	var r Robot
	_, err := fmt.Sscanf(s, "p=%d,%d v=%d,%d", &r.P.x, &r.P.y, &r.V.x, &r.V.y)
	if err != nil {
		debugf("bad format: %q", s)
		return nil, err
	}
	return &r, nil
}

func move(R []*Robot, w, h int) {
	for _, r := range R {
		nx, ny := (r.P.x+r.V.x)%w, (r.P.y+r.V.y)%h
		if nx < 0 {
			nx += w
		}
		if nx >= w {
			nx -= w
		}
		if ny < 0 {
			ny += h
		}
		r.P.x = nx
		r.P.y = ny
	}
}

const (
	NSEC    = 100
	MINBLOB = 50
)

func printRoboField(R []*Robot, w, h int) string {
	F := makeNumField[int](h, w)
	for _, r := range R {
		F[r.P.y][r.P.x]++
	}
	return printNumFieldWithSubs(F, "", map[int]string{
		0: ".",
	})
}

func maxBlob(R []*Robot) int {
	F := make(map[Point2]int)
	for _, r := range R {
		F[r.P]++
	}

	V := make(map[Point2]bool)
	var recurse func(p Point2) int
	recurse = func(p Point2) int {
		V[p] = true
		blob := F[p]
		for _, s := range STEPS8 {
			np := Point2{
				x: p.x + s[0],
				y: p.y + s[1],
			}
			if F[np] > 0 && !V[np] {
				blob += recurse(np)
			}
		}
		return blob
	}

	maxblob := 0
	for _, r := range R {
		if V[r.P] {
			continue
		}
		maxblob = max(maxblob, recurse(r.P))
	}

	return maxblob
}

func main() {
	lines := input()

	var w, h int
	_, err := fmt.Sscanf(lines[0], "size=%d,%d", &w, &h)
	noerr(err)

	R := make([]*Robot, 0, len(lines)-2)
	for i := 2; i < len(lines); i++ {
		r, err := parseRobot(lines[i])
		noerr(err)
		R = append(R, r)
	}

	debugf("size: %dx%d", w, h)
	debugf("robots: %+v", R)

	for i := 0; i < NSEC; i++ {
		move(R, w, h)
	}

	debugf("final robots: %+v", R)

	Q := make(map[int]map[int]int)
	for _, i := range []int{-1, 0, 1} {
		Q[i] = make(map[int]int)
	}
	for _, r := range R {
		Q[cmp(r.P.x, w/2)][cmp(r.P.y, h/2)]++
	}

	debugf("Q=%+v", Q)

	res1 := 1
	for _, x := range []int{-1, 1} {
		for _, y := range []int{-1, 1} {
			res1 *= Q[x][y]
		}
	}

	printf("res1: %d", res1)

	// This code makes an assumption that the pattern emerges after the first N seconds
	step := NSEC
	for {
		step++
		move(R, w, h)
		blob := maxBlob(R)
		// The code makes an assumption that the lookup pattern is an 8-connected blob of pixels
		if blob > MINBLOB {
			printf("Interesting result at step %d", step)
			fmt.Println(printRoboField(R, w, h))
			printf("res2: %d", step)
			break
		}
	}
}

func cmp[N Number](a, b N) int {
	if a < b {
		return -1
	} else if a == b {
		return 0
	}
	return 1
}
