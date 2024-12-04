package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	ALOT    = int(999999999)
	ALOT32u = uint32(4294967295)
	ALOT32  = int32(2147483647)
	ALOT64u = uint64(18446744073709551615)
	ALOT64  = int64(9223372036854775807)
)

var (
	STEPS4 = [][2]int{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}

	STEPS8 = [][2]int{
		{-1, -1},
		{-1, 0},
		{-1, 1},
		{0, -1},
		{0, 1},
		{1, -1},
		{1, 0},
		{1, 1},
	}
)

var (
	DEBUG = os.Getenv("DEBUG") == "1"
	INPUT = os.Getenv("INPUT")
)

type Number interface {
	byte | int | int32 | int64 | uint32 | uint64 | float64
}

type Integer interface {
	int | int8 | uint8 | int16 | uint16 | int32 | int64 | uint32 | uint64
}

func input() []string {
	if len(INPUT) == 0 {
		fatalf("`INPUT` var is missing")
	}

	f, err := os.Open(INPUT)
	noerr(err)
	defer f.Close()

	return readLines(f)
}

func noerr(err error) {
	if err != nil {
		panic(fmt.Sprintf("unhandled error: %s", err))
	}
}

func assert(check bool, msg string) {
	if !check {
		panic(fmt.Sprintf("assert %q failed", msg))
	}
}

func parseInt(s string) int {
	num, err := strconv.Atoi(s)
	noerr(err)
	return num
}

func readFile(in io.Reader) string {
	data, err := ioutil.ReadAll(in)
	noerr(err)
	return trim(string(data))
}

func readLines(in io.Reader) []string {
	scanner := bufio.NewScanner(in)
	lines := make([]string, 0, 1)
	for scanner.Scan() {
		lines = append(lines, trim(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scan failed: %s", err))
	}
	return lines
}

func trim(s string) string {
	return strings.TrimRight(s, "\t\n\r")
}

func parseInts(s string) []int {
	chs := strings.FieldsFunc(trim(s), func(r rune) bool {
		return r == ' ' || r == ',' || r == '\t'
	})
	nums := make([]int, 0, len(chs))
	for i := 0; i < len(chs); i++ {
		nums = append(nums, parseInt(chs[i]))
	}
	return nums
}

func numDigs(n int) int {
	d := 1
	for n >= 10 {
		d += 1
		n /= 10
	}
	return d
}

func glueNum(nn []int) int {
	N := nn[0]
	for i := 1; i < len(nn); i++ {
		for k := 0; k < numDigs(nn[i]); k++ {
			N *= 10
		}
		N += nn[i]
	}
	return N
}

func makeNumField[N Number](h, w int) [][]N {
	res := make([][]N, h)
	for i := 0; i < h; i++ {
		res[i] = make([]N, w)
	}
	return res
}

func makeIntField(h, w int) [][]int {
	return makeNumField[int](h, w)
}

func makeByteField(h, w int) [][]byte {
	return makeNumField[byte](h, w)
}

func sizeNumField[N Number](field [][]N) (int, int) {
	rows, cols := len(field), 0
	if rows > 0 {
		cols = len(field[0])
	}
	return rows, cols
}

// Deprecated: please use `sizeNumField` instead.
func sizeIntField(field [][]int) (int, int) {
	return sizeNumField(field)
}

// Deprecated: please use `sizeNumField` instead.
func sizeByteField(field [][]byte) (int, int) {
	return sizeNumField(field)
}

func copyNumField[N Number](field [][]N) [][]N {
	cp := makeNumField[N](sizeNumField(field))
	for i := 0; i < len(field); i++ {
		copy(cp[i], field[i])
	}
	return cp
}

func countNumField[N Number](field [][]N) N {
	var cnt N
	for i := 0; i < len(field); i++ {
		for j := 0; j < len(field[0]); j++ {
			cnt += field[i][j]
		}
	}
	return cnt
}

// Deprecated: please use `copyNumField` instead.
func copyIntField(field [][]int) [][]int {
	return copyNumField(field)
}

// Deprecated: please use `copyNumField` instead.
func copyByteField(field [][]byte) [][]byte {
	return copyNumField(field)
}

func printNumField[N Number](field [][]N, sep string) string {
	return printNumFieldWithSubs(field, sep, make(map[N]string))
}

// Deprecated: please use `printNumField` instead.
func printIntField(field [][]int, sep string) string {
	return printNumFieldWithSubs(field, sep, make(map[int]string))
}

// Deprecated: please use `printNumField` instead.
func printByteField(field [][]byte, sep string) string {
	return printNumFieldWithSubs(field, sep, make(map[byte]string))
}

func printNumFieldWithSubs[N Number](field [][]N, sep string, subs map[N]string) string {
	var buf bytes.Buffer
	rows, cols := sizeNumField(field)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if j > 0 {
				buf.WriteString(sep)
			}
			if sub, ok := subs[field[i][j]]; ok {
				buf.WriteString(sub)
			} else {
				buf.WriteByte('0' + byte(field[i][j]))
			}
		}
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	return buf.String()
}

func printIntFieldWithSubs(field [][]int, sep string, subs map[int]string) string {
	return printNumFieldWithSubs(field, sep, subs)
}

func printByteFieldWithSubs(field [][]byte, sep string, subs map[byte]string) string {
	return printNumFieldWithSubs(field, sep, subs)
}

func print2DMapWithSubs[T comparable](M map[Point2]T, subs map[T]string) string {
	var buf strings.Builder
	var minx, miny, maxx, maxy = ALOT, ALOT, -ALOT, -ALOT

	for p := range M {
		minx = min(minx, p.x)
		maxx = max(maxx, p.x)
		miny = min(miny, p.x)
		maxy = max(maxy, p.x)
	}

	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			p := Point2{x: x, y: y}
			if v, ok := M[p]; ok {
				if s, ok := subs[v]; ok {
					buf.WriteString(s)
				} else {
					buf.WriteByte('?')
				}
			} else {
				buf.WriteByte('.')
			}
		}
		buf.WriteRune('\n')
	}

	return buf.String()
}

func min[N Number](a, b N) N {
	if a < b {
		return a
	}
	return b
}

func max[N Number](a, b N) N {
	if a > b {
		return a
	}
	return b
}

func abs[N Number](v N) N {
	if v < 0 {
		return -v
	}
	return v
}

// functions to compute local extremums

func findLocalMin(n int, compute func(i int) int) int {
	a, b := 0, n-1
	leftix, midix, rightix := a, (a+b)/2, b
	left, mid, right := compute(leftix), compute(midix), compute(rightix)
	for rightix-leftix > 1 {
		if left <= mid && mid <= right {
			b = midix
			leftix, midix, rightix = a, (a+midix)/2, midix
			left, mid, right = compute(leftix), compute(midix), mid
		} else if left >= mid && mid >= right {
			a = midix
			leftix, midix, rightix = midix, (midix+b)/2, b
			left, mid, right = right, compute(midix), compute(rightix)
		} else {
			a = leftix
			b = rightix
			leftix, rightix = (leftix+midix)/2, (midix+rightix)/2
			left, right = compute(leftix), compute(rightix)
		}
	}
	return min(left, right)
}

func findLocalMax(n int, compute func(i int) int) int {
	return -1 * findLocalMin(n, func(i int) int {
		return -1 * compute(i)
	})
}

// slice helpers

func mapIntArr(arr []int, mapfn func(int) int) []int {
	res := make([]int, len(arr))
	for i := 0; i < len(arr); i++ {
		res[i] = mapfn(arr[i])
	}
	return res
}

func mapByteArr(arr []byte, mapfn func(byte) byte) []byte {
	res := make([]byte, len(arr))
	for i := 0; i < len(arr); i++ {
		res[i] = mapfn(arr[i])
	}
	return res
}

func reverseNumArr[N Number](arr []N) []N {
	res := make([]N, len(arr))
	for i := 0; i < len(arr); i++ {
		res[len(arr)-1-i] = arr[i]
	}
	return res
}

// Deprecated: please use `reverseNumArr` instead.
func reverseIntArr(arr []int) []int {
	return reverseNumArr(arr)
}

// Deprecated: please use `reverseNumArr` instead.
func reverseByteArr(arr []byte) []byte {
	return reverseByteArr(arr)
}

func reverseStr(s string) string {
	rs := []rune(s)
	return string(reverseNumArr(rs))
}

func grepNumArr[N Number](arr []N, grepfn func(N) bool) []N {
	res := make([]N, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		if grepfn(arr[i]) {
			res = append(res, arr[i])
		}
	}
	return res
}

// Deprecated: please use `grepNumArr` instead.
func grepIntArr(arr []int, grepfn func(int) bool) []int {
	return grepNumArr(arr, grepfn)
}

// Deprecated: please use `grepNumArr` instead.
func grepByteArr(arr []byte, grepfn func(byte) bool) []byte {
	return grepNumArr(arr, grepfn)
}

func transposeMat[N Number](mx [][]N) [][]N {
	h, w := sizeNumField(mx)
	cp := makeNumField[N](w, h)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			cp[j][i] = mx[i][j]
		}
	}
	return cp
}

func reverseMatHor[N Number](mx [][]N) [][]N {
	h, w := sizeNumField(mx)
	cp := makeNumField[N](h, w)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			cp[i][w-1-j] = mx[i][j]
		}
	}
	return cp
}

func reverseMatVer[N Number](mx [][]N) [][]N {
	h, w := sizeNumField(mx)
	cp := makeNumField[N](h, w)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			cp[h-1-i][j] = mx[i][j]
		}
	}
	return cp
}

func rotateMatLeft[N Number](mx [][]N) [][]N {
	return reverseMatVer(transposeMat(mx))
}

func rotateMatRight[N Number](mx [][]N) [][]N {
	return transposeMat(reverseMatVer(mx))
}

// logging function

func debugf(format string, v ...interface{}) {
	if DEBUG {
		log.Printf(format, v...)
	}
}

func printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func startsWith(s string, pref string) bool {
	return len(s) >= len(pref) && s[:len(pref)] == pref
}

// Data types

type BinHeap[T comparable] struct {
	items []T
	index map[T]int
	cmp   func(a, b T) bool
}

func NewBinHeap[T comparable](cmp func(a, b T) bool) *BinHeap[T] {
	return &BinHeap[T]{
		items: make([]T, 0, 1),
		index: make(map[T]int),
		cmp:   cmp,
	}
}

func (h *BinHeap[T]) Size() int {
	return len(h.items)
}

func (h *BinHeap[T]) Push(item T) {
	last := len(h.items)
	if _, ok := h.index[item]; !ok {
		h.items = append(h.items, item)
		h.index[item] = last
	}
	ptr := h.index[item]
	h.reheapAt(ptr)
}

func (h *BinHeap[T]) Pop() T {
	last := len(h.items) - 1
	h.swap(0, last)
	item := h.items[last]
	h.items = h.items[:last]
	delete(h.index, item)
	h.reheapAt(0)

	return item
}

func (h *BinHeap[T]) swap(i, j int) {
	h.index[h.items[i]], h.index[h.items[j]] = h.index[h.items[j]], h.index[h.items[i]]
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *BinHeap[T]) reheapAt(ptr int) {
	for ptr > 0 {
		parent := (ptr - 1) / 2
		if h.cmp(h.items[ptr], h.items[parent]) {
			h.swap(ptr, parent)
			ptr = parent
		} else {
			break
		}
	}

	for ptr < len(h.items) {
		ch1, ch2 := ptr*2+1, ptr*2+2
		next := ptr
		if ch1 < len(h.items) && h.cmp(h.items[ch1], h.items[next]) {
			next = ch1
		}
		if ch2 < len(h.items) && h.cmp(h.items[ch2], h.items[next]) {
			next = ch2
		}
		if next != ptr {
			h.swap(ptr, next)
			ptr = next
		} else {
			break
		}
	}
}

// 2d and 3d points

type Point2 struct {
	x, y int
}

func NewPoint2(x, y int) *Point2 {
	return &Point2{x, y}
}

func (p2 *Point2) String() string {
	return fmt.Sprintf("P2{%d, %d}", p2.x, p2.y)
}

func (p2 *Point2) Dist(other Point2) int {
	return abs(p2.x-other.x) + abs(p2.y-other.y)
}

type Point3 struct {
	x, y, z int
}

func NewPoint3(x, y, z int) *Point3 {
	return &Point3{x, y, z}
}

func (p3 *Point3) String() string {
	return fmt.Sprintf("P3{%d, %d, %d}", p3.x, p3.y, p3.z)
}

func (p3 *Point3) Dist(other Point3) int {
	return abs(p3.x-other.x) + abs(p3.y-other.y) + abs(p3.z-other.z)
}

// Ranges

type Range [2]int

func NewRange(a, b int) Range {
	assert(a <= b, "Range a must be less or equal to b")
	return [2]int{a, b}
}

func (r Range) Intersects(other Range) bool {
	return (r[0] >= other[0] && r[1] <= other[1]) || (r[1] >= other[0] && r[1] <= other[1]) || (other[0] >= r[0] && other[1] <= r[1]) || (other[1] >= r[0] && other[1] <= r[1])
}

func (r Range) Contains(v int) bool {
	return r[0] >= v && r[1] <= v
}

func (r Range) String() string {
	return fmt.Sprintf("[%d..%d]", r[0], r[1])
}

func mergeRanges(ranges []Range) []Range {
	sort.Slice(ranges, func(a, b int) bool {
		return ranges[a][0] < ranges[b][0]
	})

	ix := 0
	for ix < len(ranges)-1 {
		if ranges[ix][1] >= ranges[ix+1][0] {
			ranges[ix][1] = max(ranges[ix][1], ranges[ix+1][1])
			ranges = append(ranges[:ix+1], ranges[ix+2:]...)
		} else {
			ix++
		}
	}
	return ranges
}

// Parse functions

func peek(s string, ptr int) byte {
	return s[ptr]
}

func match(s string, ptr int, b byte) bool {
	return ptr < len(s) && peek(s, ptr) == b
}

func matchStr(s string, ptr int, lex string) bool {
	if len(lex) > len(s)-ptr {
		return false
	}
	return s[ptr:ptr+len(lex)] == lex
}

func consume(s string, ptr int, b byte) int {
	if match(s, ptr, b) {
		return ptr + 1
	}
	panic(fmt.Sprintf("consume mismatch at pos: %d around %s", ptr, s[:ptr+1]))
}

func eatWhitespace(s string, ptr int) int {
	for ptr < len(s) && isWhitespace(s, ptr) {
		ptr++
	}
	return ptr
}

func readInt(s string, ptr int) (int, int) {
	from := ptr
	if match(s, ptr, '+') || match(s, ptr, '-') {
		ptr++
	}
	for ptr < len(s) && isNumber(s[ptr]) {
		ptr++
	}
	return parseInt(s[from:ptr]), ptr
}

func readFloat64(s string, ptr int) (float64, int) {
	from := ptr
	comma := 0
	if match(s, ptr, '+') || match(s, ptr, '-') {
		ptr++
	}
	for ptr < len(s) && isNumber(s[ptr]) || (comma < 1 && match(s, ptr, '.')) {
		if match(s, ptr, '.') {
			comma++
		}
		ptr++
	}
	f, err := strconv.ParseFloat(s[from:ptr], 64)
	noerr(err)
	return f, ptr
}

func readStr(s string, ptr int, lex string) (string, int) {
	off := 0
	for (ptr+off) < len(s) && off < len(lex) {
		if s[ptr+off] != lex[off] {
			panic(fmt.Sprintf("readStr mismatch at pos %d around %s", ptr+off, s[:ptr+off+1]))
		}
		off++
	}
	return s[ptr : ptr+off], ptr + off
}

func readWord(s string, ptr int) (string, int) {
	from := ptr
	for ptr < len(s) && (isAlpha(s[ptr]) || isNumber(s[ptr])) {
		ptr++
	}
	return s[from:ptr], ptr
}

func isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlpha(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isAlphaNum(b byte) bool {
	return isAlpha(b) || isNumber(b)
}

func isWhitespace(s string, ptr int) bool {
	return s[ptr] == ' ' || s[ptr] == '\t' || s[ptr] == '\r'
}

func computeSubsetsOfLenRange[Item any](items []Item, low, high int) [][]Item {
	res := make([][]Item, 0, 1)
	for n := low; n <= high; n++ {
		res = append(res, computeSubsetsOfLen(items, n)...)
	}
	return res
}

func computeSubsetsOfLen[Item any](items []Item, n int) [][]Item {
	var recurse func(cur, n int) [][]Item
	recurse = func(cur, n int) [][]Item {
		if cur >= len(items) {
			return nil
		}
		if n == 0 {
			return [][]Item{{}}
		}
		res := [][]Item{}
		for _, sub := range recurse(cur+1, n-1) {
			res = append(res, append([]Item{items[cur]}, sub...))
		}
		for _, sub := range recurse(cur+1, n) {
			res = append(res, sub)
		}
		return res
	}

	return recurse(0, n)
}

func execScanfInstr(instr string, input string, iptr int, arg any) int {
	//TODO: check the pointer type
	switch instr {
	case "byte":
		*(arg.(*byte)) = input[iptr]
		iptr++
	case "int":
		var intv int
		intv, iptr = readInt(input, iptr)
		*(arg.(*int)) = intv
	case "float":
		var floatv float64
		floatv, iptr = readFloat64(input, iptr)
		*(arg.(*float64)) = floatv
	case "string":
		var strv string
		strv, iptr = readWord(input, iptr)
		*(arg.(*string)) = strv
	case "[]int":
		arr := make([]int, 0, 1)
		var num int
		for iptr < len(input) {
			num, iptr = readInt(input, iptr)
			arr = append(arr, num)
			if !match(input, iptr, ',') {
				break
			}
			iptr++
			iptr = eatWhitespace(input, iptr)
		}
		*(arg.(*[]int)) = arr
	case "[]float":
		arr := make([]float64, 0, 1)
		var num float64
		for iptr < len(input) {
			num, iptr = readFloat64(input, iptr)
			arr = append(arr, num)
			if !match(input, iptr, ',') {
				break
			}
			iptr++
			iptr = eatWhitespace(input, iptr)
		}
		*(arg.(*[]float64)) = arr
	case "[]string":
		arr := make([]string, 0, 1)
		var str string
		for iptr < len(input) {
			str, iptr = readWord(input, iptr)
			arr = append(arr, str)
			if !match(input, iptr, ',') {
				break
			}
			iptr++
			iptr = eatWhitespace(input, iptr)
		}
		*(arg.(*[]string)) = arr
	default:
		fatalf("unsupported instr: %s", instr)
	}
	return iptr
}

func Scanf(input string, tmpl string, args ...any) {
	iptr, tptr := 0, 0
	argix := 0
	for iptr < len(input) && tptr < len(tmpl) {
		for tptr < len(tmpl) && iptr < len(input) && tmpl[tptr] != '{' {
			if input[iptr] != tmpl[tptr] {
				fatalf("input and template mismatch around pos %d: %s", iptr, input[:iptr+1])
			}
			tptr++
			iptr++
		}
		if tptr < len(tmpl) && tmpl[tptr] == '{' {
			tptr++
			from := tptr
			for tptr < len(tmpl) && tmpl[tptr] != '}' {
				tptr++
			}
			if tmpl[tptr] != '}' {
				fatalf("template is missing a terminating curly brace: %s", tmpl[:tptr])
			}
			instr := tmpl[from:tptr]
			tptr++
			iptr = execScanfInstr(instr, input, iptr, args[argix])
			argix++
		}
	}
}

const (
	EOW  rune = 0
	ROOT      = 255
)

type TrieNode struct {
	r     rune
	nodes map[rune]*TrieNode
}

func NewTrieNode(r rune) *TrieNode {
	return &TrieNode{
		r:     r,
		nodes: make(map[rune]*TrieNode),
	}
}

type Trie struct {
	head *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		head: NewTrieNode(ROOT),
	}
}

func (t *Trie) Add(word string) {
	ptr := t.head
	for _, r := range word {
		if _, ok := ptr.nodes[r]; !ok {
			ptr.nodes[r] = NewTrieNode(r)
		}
		ptr = ptr.nodes[r]
	}
	ptr.nodes[EOW] = NewTrieNode(EOW)
}

func (t *Trie) AddAll(words ...string) {
	for _, word := range words {
		t.Add(word)
	}
}

type mtch struct {
	ix   int
	tptr *TrieNode
}

func FirstIndexOfAny(s string, words ...string) (string, int) {
	t := NewTrie()
	t.AddAll(words...)
	ms := make([]mtch, 0, 1)
	for ix, r := range s {
		ms = append(ms, mtch{ix, t.head})
		newms := make([]mtch, 0, 1)
		for _, m := range ms {
			if next, ok := m.tptr.nodes[r]; ok {
				if _, ok := next.nodes[EOW]; ok {
					return s[m.ix : ix+1], m.ix
				}
				newms = append(newms, mtch{m.ix, next})
			}
		}
		ms = newms
	}
	return "", -1
}

func LastIndexOfAny(s string, words ...string) (string, int) {
	sinv := reverseStr(s)
	winv := make([]string, 0, len(words))
	for _, word := range words {
		winv = append(winv, reverseStr(word))
	}
	sres, ix := FirstIndexOfAny(sinv, winv...)
	return reverseStr(sres), len(s) - ix - 1
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(a, b int) int {
	return a * b / GCD(a, b)
}

func popcount[N Integer](n N) int {
	cnt := 0
	for n > 0 {
		n &= (n - 1)
		cnt++
	}
	return cnt
}
