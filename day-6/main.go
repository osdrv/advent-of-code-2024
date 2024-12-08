package main

import (
	"sort"
)

const (
	DOWN  = 'v'
	UP    = '^'
	LEFT  = '<'
	RIGHT = '>'
)

var TURN = map[byte]byte{
	UP:    RIGHT,
	RIGHT: DOWN,
	DOWN:  LEFT,
	LEFT:  UP,
}

var MOVE = map[byte][2]int{
	UP:    {-1, 0},
	RIGHT: {0, 1},
	DOWN:  {1, 0},
	LEFT:  {0, -1},
}

const (
	OBSTACLE = '#'
	EMPTY    = '.'
)

func traverse(F [][]byte, start [2]int, sdir byte) (int, bool) {
	v2 := make(map[uint16]struct{})
	v3 := make(map[uint32]struct{})

	pos := start
	dir := sdir
	var p2 uint16
	var p3 uint32
	for pos[0] >= 0 && pos[0] < len(F) && pos[1] >= 0 && pos[1] < len(F[pos[0]]) {
		p2 = uint16(pos[0])<<8 + uint16(pos[1])
		p3 = uint32(pos[0])<<16 + uint32(pos[1])<<8 + uint32(dir)
		v2[p2] = struct{}{}
		if _, ok := v3[p3]; ok {
			return -1, false
		}
		v3[p3] = struct{}{}
		move := MOVE[dir]
		ni, nj := pos[0]+move[0], pos[1]+move[1]
		if ni < 0 || ni >= len(F) || nj < 0 || nj >= len(F[ni]) {
			break
		}
		if F[ni][nj] == OBSTACLE {
			dir = TURN[dir]
			continue
		}
		pos = [2]int{ni, nj}
	}

	return len(v2), true
}

func prevGt[N Number](nums []N, n N) int {
	if len(nums) == 0 || n < nums[0] {
		return -1
	}
	if n > nums[len(nums)-1] {
		return len(nums) - 1
	}
	return sort.Search(len(nums), func(i int) bool {
		return nums[i] > n
	}) - 1
}

func nextLw[N Number](nums []N, n N) int {
	if len(nums) == 0 || n > nums[len(nums)-1] {
		return -1
	}
	if n < nums[0] {
		return 0
	}
	return sort.Search(len(nums), func(i int) bool {
		return nums[i] >= n
	})
}

func insert(a []int, x int) []int {
	res := make([]int, len(a)+1)
	ix := sort.Search(len(a), func(i int) bool {
		return a[i] > x
	})
	copy(res[:ix], a[:ix])
	res[ix] = x
	copy(res[ix+1:], a[ix:])
	return res
}

func hasLoop(vert, hor [][]int, start [2]int, sdir byte) bool {
	dir := sdir
	pos := start

	visited := make(map[uint64]bool)

	for {
		i, j := pos[0], pos[1]

		vix := uint64(i)<<32 + uint64(j)<<16 + uint64(dir)
		if _, ok := visited[vix]; ok {
			return true
		}
		visited[vix] = true

		ni, nj := i, j
		switch dir {
		case UP:
			nix := prevGt(vert[j], i)
			if nix == -1 {
				goto NOLOOP
			}
			ni = vert[j][nix] + 1
		case DOWN:
			nix := nextLw(vert[j], i)
			if nix == -1 {
				goto NOLOOP
			}
			ni = vert[j][nix] - 1
		case LEFT:
			nix := prevGt(hor[i], j)
			if nix == -1 {
				goto NOLOOP
			}
			nj = hor[i][nix] + 1
		case RIGHT:
			nix := nextLw(hor[i], j)
			if nix == -1 {
				goto NOLOOP
			}
			nj = hor[i][nix] - 1
		default:
			panic("wtf")
		}
		pos = [2]int{ni, nj}
		dir = TURN[dir]
	}
NOLOOP:
	return false
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	start := [2]int{-1, -1}
	dir := byte(0)
	obsts := make([][2]int, 0)

	hor := make([][]int, len(lines))
	vert := make([][]int, len(lines[0]))

	for i, line := range lines {
		F = append(F, []byte(line))
		for j := 0; j < len(F[i]); j++ {
			if _, ok := MOVE[F[i][j]]; ok {
				start = [2]int{i, j}
				dir = F[i][j]
				F[i][j] = EMPTY
			}
			if F[i][j] == OBSTACLE {
				obsts = append(obsts, [2]int{i, j})
				hor[i] = append(hor[i], j)
				vert[j] = append(vert[j], i)
			}
		}
	}

	steps1, _ := traverse(F, start, dir)
	printf("steps1: %d", steps1)

	obst := 0
	for i := 0; i < len(F); i++ {
		for j := 0; j < len(F[i]); j++ {
			if F[i][j] == EMPTY {
				F[i][j] = OBSTACLE
				vertbak := vert[j]
				horbak := hor[i]

				vert[j] = insert(vert[j], i)
				hor[i] = insert(hor[i], j)

				if hasLoop(vert, hor, start, dir) {
					debugf("extra obst at i: %d, j: %d creates a loop", i, j)
					obst++
				}
				F[i][j] = EMPTY
				vert[j] = vertbak
				hor[i] = horbak
			}
		}
	}
	printf("obst2: %d", obst)
}
