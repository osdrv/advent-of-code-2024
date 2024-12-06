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
	EXTRA    = 'O'
)

const LOOP_BREAKER = 4 // magic number, I started with 10 and reduced it down to 4.
// It might be different for other inputs.

func traverse(F [][]byte, start [2]int, sdir byte, loopcnt int) (int, bool) {
	visited := make(map[[2]int]int)
	maxv := 0

	pos := start
	dir := sdir
	for pos[0] >= 0 && pos[0] < len(F) && pos[1] >= 0 && pos[1] < len(F[pos[0]]) {
		visited[pos]++
		if visited[pos] > loopcnt {
			// I guess, this is a loop
			return -1, false
		}
		maxv = max(maxv, visited[pos])
		move := MOVE[dir]
		ni, nj := pos[0]+move[0], pos[1]+move[1]
		if ni < 0 || ni >= len(F) || nj < 0 || nj >= len(F[ni]) {
			break
		}
		if F[ni][nj] == OBSTACLE {
			dir = TURN[dir]
			continue
		}
		F[ni][nj] = 'X'
		pos = [2]int{ni, nj}
	}

	return len(visited), true
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

	FCP := copyNumField(F)
	steps1, _ := traverse(FCP, start, dir, LOOP_BREAKER)
	printf("steps1: %d", steps1)

	obst := 0
	for i := 0; i < len(F); i++ {
		for j := 0; j < len(F[i]); j++ {
			FCP := copyNumField(F)
			if FCP[i][j] == EMPTY {
				FCP[i][j] = OBSTACLE

				if _, ok := traverse(FCP, start, dir, LOOP_BREAKER); !ok {
					// found a loop
					obst++
					FCP[i][j] = EXTRA
				}
			}
		}
	}
	printf("obst: %d", obst)
}
