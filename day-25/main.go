package main

func parseKey(lines []string) []int {
	res := make([]int, len(lines[0]))
	for j := 0; j < len(lines[0]); j++ {
		for i := 0; i < len(lines)-1; i++ {
			if lines[len(lines)-2-i][j] == '.' {
				res[j] = i
				break
			}
		}
	}
	return res
}

func parseLock(lines []string) []int {
	res := make([]int, len(lines[0]))
	for j := 0; j < len(lines[0]); j++ {
		for i := 0; i < len(lines)-1; i++ {
			if lines[i+1][j] == '.' {
				res[j] = i
				break
			}
		}
	}
	return res
}

func wouldFit(height int, lock, key []int) bool {
	for i := 0; i < len(lock); i++ {
		if lock[i]+key[i] > height {
			return false
		}
	}
	return true
}

func main() {
	lines := input()

	locks := make([][]int, 0, 1)
	keys := make([][]int, 0, 1)

	prev := 0
	ix := 0
	height := 0
	for ix <= len(lines) {
		if ix == len(lines) || len(lines[ix]) == 0 {
			if lines[prev][0] == '#' {
				// lock
				locks = append(locks, parseLock(lines[prev:ix]))
			} else if lines[ix-1][0] == '#' {
				// key
				keys = append(keys, parseKey(lines[prev:ix]))
			}
			height = ix - prev
			ix++
			prev = ix
			continue
		}
		ix++
	}
	debugf("locks: %+v", locks)
	debugf("keys: %+v", keys)
	debugf("height: %d", height)

	res1 := 0
	for _, lock := range locks {
		for _, key := range keys {
			// height-2 because height is the total height of the tile and
			// the upper and the lower edges are eaten
			if wouldFit(height-2, lock, key) {
				debugf("lock %+v and key %+v would fit", lock, key)
				res1++
			}
		}
	}

	printf("res1: %d", res1)
}
