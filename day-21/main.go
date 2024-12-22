package main

import (
	"sort"
	"strconv"
	"strings"
)

var (
	LEX    = make(map[string]int)
	LEXCNT = 0
	LEXIX  = make(map[int]string)

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

	MOVES = map[[2]int]byte{
		{0, 1}:  '>',
		{0, -1}: '<',
		{1, 0}:  'v',
		{-1, 0}: '^',
		{0, 0}:  'A',
	}
)

func Lookup(n int) string {
	return LEXIX[n]
}

func Lexem(s string) int {
	if _, ok := LEX[s]; !ok {
		LEXIX[LEXCNT] = s
		LEX[s] = LEXCNT
		LEXCNT++
	}
	return LEX[s]
}

type Path struct {
	Segms  map[int]uint64
	Len    uint64
	_str   string
	_score uint64
}

func NewPath() *Path {
	return &Path{
		Segms:  make(map[int]uint64),
		Len:    0,
		_str:   "",
		_score: 0,
	}
}

func (p *Path) Copy() *Path {
	cp := NewPath()
	cp.Len = p.Len
	for n, cnt := range p.Segms {
		cp.Segms[n] = cnt
	}
	cp._str = p._str
	cp._score = p._score
	return cp
}

func (p *Path) AddSegm(segm string, cnt uint64) {
	p._str = ""
	p._score = 0
	n := Lexem(segm)
	p.Segms[n] += cnt
	p.Len += uint64(len(segm)) * cnt
}

func (p *Path) String() string {
	if len(p._str) > 0 {
		return p._str
	}
	nums := make([]int, 0, len(p.Segms))
	for segm := range p.Segms {
		nums = append(nums, segm)
	}
	sort.Slice(nums, func(i, j int) bool {
		return Lookup(nums[i]) < Lookup(nums[j])
	})
	var b strings.Builder
	b.WriteString("Path[Segms:{")
	for i, num := range nums {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(Lookup(num))
		b.WriteByte(':')
		b.WriteString(strconv.FormatUint(p.Segms[num], 10))
	}
	b.WriteString("} Len:")
	b.WriteString(strconv.FormatUint(p.Len, 10))
	b.WriteByte(']')

	p._str = b.String()

	return p._str
}

// The score is based on the observation that paths start deviate in length after the second path
func (p *Path) Score() uint64 {
	if p._score > 0 {
		return p._score
	}
	score := ALOT64u
	for n, cnt := range p.Segms {
		subpp := traverse(DIRECTIONAL, Point2{x: 2, y: 0}, Lookup(n))
		minsubsubscore := ALOT64u
		for _, subp := range subpp {
			subsubscore := uint64(0)
			for subn, subcnt := range subp.Segms {
				subsubpp := traverse(DIRECTIONAL, Point2{x: 2, y: 0}, Lookup(subn))
				subsubscore += subsubpp[0].Len * subcnt
			}
			minsubsubscore = min(minsubsubscore, subsubscore)
		}
		score += minsubsubscore * cnt
	}
	p._score = score
	return score
}

func ParsePath(s string) *Path {
	path := NewPath()
	ix := 0
	pix := 0
	for ix < len(s) {
		if s[ix] == 'A' {
			path.AddSegm(s[pix:ix+1], 1)
			ix++
			pix = ix
		} else {
			ix++
		}
	}
	if ix != pix {
		path.AddSegm(s[pix:ix+1], 1)
	}
	return path
}

type QItem struct {
	P    Point2
	Ix   int
	Path []byte
}

var (
	MEMO = make(map[string][]*Path)
)

func traverse(F [][]byte, start Point2, path string) []*Path {
	if prev, ok := MEMO[path]; ok {
		return prev
	}
	Q := make([]QItem, 0, 1)
	Q = append(Q, QItem{
		P:    start,
		Ix:   0,
		Path: []byte{},
	})

	V := make(map[Point3]int)

	opts := make([]string, 0, 1)
	minlen := ALOT
	var H QItem
	for len(Q) > 0 {
		H, Q = Q[0], Q[1:]
		p3 := Point3{x: H.P.x, y: H.P.y, z: H.Ix}
		if prev, ok := V[p3]; ok && prev < len(H.Path) {
			continue
		}
		V[p3] = len(H.Path)
		if F[H.P.y][H.P.x] == path[H.Ix] {
			H.Path = append(H.Path, 'A')
			H.Ix++
			if H.Ix >= len(path) {
				if minlen >= len(H.Path) {
					if minlen > len(H.Path) {
						opts = make([]string, 0, 1)
						minlen = len(H.Path)
					}
					opts = append(opts, string(H.Path))
				}
			} else {
				Q = append(Q, H)
			}
			continue
		}

		for _, step := range STEPS4 {
			ny, nx := H.P.y+step[0], H.P.x+step[1]
			if ny < 0 || ny >= len(F) || nx < 0 || nx >= len(F[0]) || F[ny][nx] == 0 {
				continue
			}
			p3 = Point3{x: nx, y: ny, z: H.Ix}
			newpath := make([]byte, len(H.Path)+1)
			copy(newpath, H.Path)
			newpath[len(newpath)-1] = MOVES[step]
			if prev, ok := V[p3]; ok && prev < len(newpath) {
				continue
			}
			Q = append(Q, QItem{
				P:    Point2{x: nx, y: ny},
				Ix:   H.Ix,
				Path: newpath,
			})
		}
	}

	res := make([]*Path, 0, len(opts))
	for _, opt := range opts {
		res = append(res, ParsePath(opt))
	}

	MEMO[path] = res

	return res
}

var (
	PMEMO = make(map[string][]*Path)
)

func traversePath(F [][]byte, start Point2, path *Path) []*Path {
	if prev, ok := PMEMO[path.String()]; ok {
		debugf("memo hit!")
		return prev
	}

	newpaths := []*Path{NewPath()}

	for n, cnt := range path.Segms {
		segm := Lookup(n)
		opts := traverse(F, start, segm)
		newnewpaths := make([]*Path, 0, len(newpaths)*len(opts))

		for _, opt := range opts {
			for _, prevpath := range newpaths {
				newpath := prevpath.Copy()
				for newn, newcnt := range opt.Segms {
					newpath.AddSegm(Lookup(newn), newcnt*cnt)
				}
				newnewpaths = append(newnewpaths, newpath)
			}
		}

		newpaths = newnewpaths
	}

	res := make([]*Path, 0, 1)
	minlen := ALOT64u
	isSeen := make(map[string]bool)
	for _, newpath := range newpaths {
		if newpath.Len <= minlen {
			if newpath.Len < minlen {
				minlen = newpath.Len
				res = make([]*Path, 0, 1)
				isSeen = make(map[string]bool)
			}
			if str := newpath.String(); !isSeen[str] {
				isSeen[str] = true
				res = append(res, newpath)
			}
		}
	}

	PMEMO[path.String()] = res

	return res
}

func shortestUniqPaths(F [][]byte, start Point2, paths []*Path) []*Path {
	minlen := ALOT64u
	minpaths := make([]*Path, 0, 1)
	isSeen := make(map[string]bool)
	for _, path := range paths {
		newpaths := traversePath(F, start, path)

		for _, newpath := range newpaths {
			if minlen >= newpath.Len {
				if minlen > newpath.Len {
					minlen = newpath.Len
					minpaths = make([]*Path, 0, 1)
					isSeen = make(map[string]bool)
				}
				if str := newpath.String(); !isSeen[str] {
					isSeen[str] = true
					minpaths = append(minpaths, newpath)
				}
			}
		}
	}

	//return minpaths

	res := make([]*Path, 0, 1)
	minscore := ALOT64u
	for _, minpath := range minpaths {
		if score := minpath.Score(); score <= minscore {
			if score < minscore {
				res = make([]*Path, 0, 1)
				minscore = score
			}
			res = append(res, minpath)
		}
	}
	return res
}

func main() {
	lines := input()
	debugf("file data: %+v", lines)

	debugf("A=%+v", ParsePath("A"))
	debugf("AA=%+v", ParsePath("AA"))
	debugf("<Av<A=%+v", ParsePath("<Av<A"))

	res1 := uint64(0)
	for _, line := range lines {
		n := uint64(parseInt(line[:len(line)-1]))
		debugf("n: %d", n)
		paths := []*Path{ParsePath(line)}
		paths = shortestUniqPaths(NUMERIC, Point2{x: 2, y: 3}, paths)
		for i := 0; i < 25; i++ {
			debugf("i: %d", i)
			debugf("path len: %d", len(paths))
			paths = shortestUniqPaths(DIRECTIONAL, Point2{x: 2, y: 0}, paths)
		}
		pathlen := paths[0].Len
		debugf("%d * %d = %d", n, pathlen, n*pathlen)
		res1 += n * pathlen
	}

	printf("res1=%d", res1)
}
