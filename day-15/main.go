package main

const (
	WALL   = '#'
	ROBO   = '@'
	BOX    = 'O'
	EMPTY  = '.'
	LEFTB  = '['
	RIGHTB = ']'
)

var (
	MV = map[byte][2]int{
		'>': {0, 1},
		'v': {1, 0},
		'<': {0, -1},
		'^': {-1, 0},
	}
)

func moveBoxes(F [][]byte, start Point2, moves []byte) {
	p := start
	for i := 0; i < len(moves); i++ {
		m := moves[i]
		MQ := make([]Point2, 0, 1)
		MQ = append(MQ, p)
		OBJS := make(map[Point2]bool)
		willmove := true
		var head Point2
	NEXT:
		for len(MQ) > 0 {
			head, MQ = MQ[0], MQ[1:]
			OBJS[head] = true
			nx := head.x + MV[m][1]
			ny := head.y + MV[m][0]
			switch F[ny][nx] {
			case EMPTY:
				continue NEXT
			case WALL:
				willmove = false
				break NEXT
			case BOX:
				MQ = append(MQ, Point2{x: nx, y: ny})
			case LEFTB:
				if MV[m][1] != 0 {
					OBJS[Point2{x: nx, y: ny}] = true
				} else {
					MQ = append(MQ, Point2{x: nx, y: ny})
				}
				MQ = append(MQ, Point2{x: nx + 1, y: ny})
			case RIGHTB:
				if MV[m][1] != 0 {
					OBJS[Point2{x: nx, y: ny}] = true
				} else {
					MQ = append(MQ, Point2{x: nx, y: ny})
				}
				MQ = append(MQ, Point2{x: nx - 1, y: ny})
			}
		}

		if willmove {
			UPD := make(map[Point2]byte, len(OBJS))
			for obj := range OBJS {
				np := Point2{
					x: obj.x + MV[m][1],
					y: obj.y + MV[m][0],
				}
				UPD[np] = F[obj.y][obj.x]
				if _, ok := UPD[obj]; !ok {
					UPD[obj] = EMPTY
				}
			}
			for np, nv := range UPD {
				F[np.y][np.x] = nv
			}

			p.x += MV[m][1]
			p.y += MV[m][0]
		}
	}
}

func computeGPS(F [][]byte) int {
	sum := 0
	for i := 0; i < len(F); i++ {
		for j := 0; j < len(F[i]); j++ {
			if F[i][j] == BOX || F[i][j] == LEFTB {
				sum += 100*i + j
			}
		}
	}
	return sum
}

func main() {
	lines := input()
	F := make([][]byte, 0, len(lines))
	F2 := make([][]byte, 0, len(lines))
	var start Point2
	i := 0
	for len(lines[i]) > 0 {
		F = append(F, []byte(lines[i]))
		f2 := make([]byte, len(lines[i])*2)
		for j := 0; j < len(lines[i]); j++ {
			switch lines[i][j] {
			case ROBO:
				f2[2*j] = ROBO
				f2[2*j+1] = EMPTY
				start.y = i
				start.x = j
			case BOX:
				f2[2*j] = LEFTB
				f2[2*j+1] = RIGHTB
			default:
				f2[2*j] = lines[i][j]
				f2[2*j+1] = lines[i][j]
			}
		}
		i++
		F2 = append(F2, f2)
	}

	i++
	M := make([]byte, 0, (len(lines)-i)*len(lines[i]))
	for i < len(lines) {
		M = append(M, lines[i]...)
		i++
	}

	debugf("start: %+v", start)

	moveBoxes(F, start, M)
	res1 := computeGPS(F)
	printf("res1: %d", res1)

	moveBoxes(F2, Point2{x: start.x * 2, y: start.y}, M)
	res2 := computeGPS(F2)
	printf("res2: %d", res2)
}
