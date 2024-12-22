package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	NUMERIC = [][]byte{
		{'7', '8', '9'},
		{'4', '5', '6'},
		{'1', '2', '3'},
		{0, '0', 'A'},
	}

	DIRECTIONAL = [][]byte{
		{0, '^', 'A'},
		{'<', 'v', '>'},
	}

	PATH_MOVE = map[Point2]byte{
		{x: 1, y: 0}:  '>',
		{x: -1, y: 0}: '<',
		{x: 0, y: 1}:  'v',
		{x: 0, y: -1}: '^',
	}
)

func parseNum(s string) int {
	n, err := strconv.Atoi(s[:len(s)-1])
	noerr(err)
	return n
}

func find(F [][]byte, B byte) Point2 {
	for i := 0; i < len(F); i++ {
		for j := 0; j < len(F[i]); j++ {
			if F[i][j] == B {
				return Point2{x: j, y: i}
			}
		}
	}
	panic("wtf")
}

type QItem struct {
	P    Point2
	D    int
	IX   int
	Path []byte
}

func traverse(F [][]byte, start Point2, path string) []string {
	Q := make([]QItem, 0, 1)
	Q = append(Q, QItem{
		P:    start,
		D:    0,
		IX:   0,
		Path: []byte{},
	})

	mindist := ALOT
	minpaths := make([]string, 0, 1)
	var H QItem
	V := make(map[Point3]int)
	for len(Q) > 0 {
		H, Q = Q[0], Q[1:]
		p3 := Point3{x: H.P.x, y: H.P.y, z: H.IX}
		if prev, ok := V[p3]; ok && prev < H.D {
			continue
		}
		V[p3] = H.D
		if F[H.P.y][H.P.x] == path[H.IX] {
			H.IX++
			H.Path = append(H.Path, 'A')
			if H.IX >= len(path) {
				if mindist >= H.D {
					if mindist > H.D {
						minpaths = make([]string, 0, 1)
					}
					mindist = min(mindist, H.D)
					minpaths = append(minpaths, string(H.Path))
				}
				continue
			}
			if path[H.IX-1] == path[H.IX] {
				Q = append(Q, H)
				continue
			}
		}

		for _, s := range STEPS4 {
			nx, ny := H.P.x+s[0], H.P.y+s[1]
			if nx < 0 || nx >= len(F[0]) || ny < 0 || ny >= len(F) || F[ny][nx] == 0 {
				continue
			}
			p3 = Point3{x: nx, y: ny, z: H.IX}
			d := H.D + 1
			if prev, ok := V[p3]; ok && prev < d {
				continue
			}
			path := make([]byte, len(H.Path)+1)
			copy(path, H.Path)
			path[len(path)-1] = PATH_MOVE[Point2{x: s[0], y: s[1]}]
			Q = append(Q, QItem{
				P:    Point2{x: nx, y: ny},
				D:    d,
				IX:   H.IX,
				Path: path,
			})
		}
	}

	return minpaths
}

type Cost struct {
	S string
	C int
}

type Path struct {
	Segms map[string]int
}

func NewPath() *Path {
	return &Path{Segms: make(map[string]int)}
}

func (p *Path) Len() int {
	res := 0
	for segm, cnt := range p.Segms {
		res += len(segm) * cnt
	}
	return res
}

func (p *Path) String() string {
	return fmt.Sprintf("Path{Segms:%+v, Len: %d}", p.Segms, p.Len())
}

func (p *Path) Copy() *Path {
	scp := make(map[string]int)
	for seg, cnt := range p.Segms {
		scp[seg] = cnt
	}
	return &Path{
		Segms: scp,
	}
}

func (p *Path) Signature() string {
	keys := make([]string, 0, len(p.Segms))
	for segm := range p.Segms {
		keys = append(keys, segm)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString("Path{")
	for i, segm := range keys {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(segm)
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(p.Segms[segm]))
	}
	b.WriteByte('}')
	return b.String()
}

func ParsePath(s string) *Path {
	p := &Path{
		Segms: make(map[string]int),
	}
	pix, ix := 0, 0
	for ix < len(s) {
		if s[ix] == 'A' {
			//if pix != ix {
			p.Segms[s[pix:ix+1]]++
			//}
			ix++
			pix = ix
		} else {
			ix++
		}
	}
	if pix != ix {
		p.Segms[s[pix:ix]]++
	}
	return p
}

func findShortest(F [][]byte, MEMO map[string][]string, start Point2, path *Path) []*Path {
	cands := make([]*Path, 0, 1)
	cands = append(cands, NewPath())
	uniq := make(map[string]bool)
	for segm, cnt := range path.Segms {
		newcands := make([]*Path, 0, 1)
		var alts []string
		if prev, ok := MEMO[segm]; ok {
			alts = prev
		} else {
			alts = traverse(F, start, segm)
			MEMO[segm] = alts
		}
		minlen := ALOT
		for _, altpath := range alts {
			altp := ParsePath(altpath)
			for _, cand := range cands {
				candcp := cand.Copy()
				for newsegm, newcnt := range altp.Segms {
					candcp.Segms[newsegm] += cnt * newcnt
				}
				newcands = append(newcands, candcp)
			}
		}
		cands = make([]*Path, 0, 1)
		for _, cand := range newcands {
			if l := cand.Len(); l <= minlen {
				if l < minlen {
					minlen = l
					cands = make([]*Path, 0, 1)
				}
				if sign := cand.Signature(); !uniq[sign] {
					uniq[sign] = true
					cands = append(cands, cand)
				}
			}
		}
	}

	return cands
}

func minUniqPaths(paths []*Path) []*Path {
	res := make([]*Path, 0, 1)
	minlen := ALOT
	cmap := make(map[string]*Path)
	for _, cand := range paths {
		if cl := cand.Len(); cl <= minlen {
			if cl < minlen {
				minlen = cl
				cmap = make(map[string]*Path)
			}
			cmap[cand.Signature()] = cand
		}
	}
	for _, cand := range cmap {
		res = append(res, cand)
	}
	return res
}

func main() {
	lines := input()
	res1 := 0
	MEMO := make(map[string][]string)
	for _, line := range lines {
		paths := findShortest(NUMERIC, MEMO, Point2{x: 2, y: 3}, ParsePath(line))
		//debugf("paths: %+v", paths)
		for i := 0; i < 2; i++ {
			newpaths := make([]*Path, 0, 1)
			for _, path := range paths {
				newpaths = append(newpaths, findShortest(DIRECTIONAL, MEMO, Point2{x: 2, y: 0}, path)...)
			}
			paths = minUniqPaths(newpaths)
			//debugf("paths2: %+v", paths)
			debugf("len(paths2)=%d", len(paths))
		}

		n := parseNum(line)
		d := paths[0].Len()

		debugf("%d * %d = %d", n, d, n*d)
		res1 += n * d
	}
	printf("res1=%d", res1)
}
