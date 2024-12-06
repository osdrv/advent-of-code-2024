package main

var TURN = map[byte]byte{
	'^': '>',
	'>': 'v',
	'v': '<',
	'<': '^',
}

var MOVE = map[byte][2]int{
	'^': {-1, 0},
	'>': {0, 1},
	'v': {1, 0},
	'<': {0, -1},
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

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	start := [2]int{-1, -1}
	dir := byte(0)
	for i, line := range lines {
		F = append(F, []byte(line))
		for j := 0; j < len(F[i]); j++ {
			if _, ok := MOVE[F[i][j]]; ok {
				start = [2]int{i, j}
				dir = F[i][j]
				F[i][j] = EMPTY
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
				if _, ok := traverse(F, start, dir); !ok {
					// found a loop
					obst++
				}
				F[i][j] = EMPTY
			}
		}
	}
	printf("obst: %d", obst)
}
