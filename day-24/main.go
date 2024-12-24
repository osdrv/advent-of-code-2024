package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

func parseWire(s string) (string, int) {
	ss := strings.SplitN(s, ": ", 2)
	v := parseInt(ss[1])
	return ss[0], v
}

func parseGate(s string) (string, string, string, string) {
	var in1, typ, in2, out string
	_, err := fmt.Sscanf(s, "%s %s %s -> %s", &in1, &typ, &in2, &out)
	noerr(err)
	return in1, typ, in2, out
}

func getXYZ(state map[string]bool) (x uint64, y uint64, z uint64) {
	for k, v := range state {
		if v && (k[0] == 'x' || k[0] == 'y' || k[0] == 'z') {
			n := parseInt(k[1:])
			if k[0] == 'x' {
				x |= 1 << n
			} else if k[0] == 'y' {
				y |= 1 << n
			} else if k[0] == 'z' {
				z |= 1 << n
			}
		}
	}
	return
}

func compute(init map[string]bool, gates [][4]string) map[string]bool {
	state := make(map[string]bool)
	for k, v := range init {
		state[k] = v
	}

	rem := make([][4]string, len(gates))
	copy(rem, gates)

	for len(rem) > 0 {
		newrem := make([][4]string, 0, 1)
		for i := 0; i < len(rem); i++ {
			gate := rem[i]
			in1, ok1 := state[gate[0]]
			in2, ok2 := state[gate[2]]
			if !ok1 || !ok2 {
				newrem = append(newrem, gate)
				continue
			}
			var val bool
			switch gate[1] {
			case "AND":
				val = in1 && in2
			case "OR":
				val = in1 || in2
			case "XOR":
				val = in1 != in2
			}
			state[gate[3]] = val
		}
		rem = newrem
	}

	return state
}

type GraphNode struct {
	Name string
	Type string
	Ins  []*GraphNode
	Out  *GraphNode
}

func NewGraphNode(name string, typ string) *GraphNode {
	return &GraphNode{
		Name: name,
		Type: typ,
		Ins:  make([]*GraphNode, 0, 1),
	}
}

const (
	AND = "AND"
	OR  = "OR"
	XOR = "XOR"
)

// This function only works for the input, not for test data because it expects the adder structure
func fixGraph(gates [][4]string, subs map[string]string) {
	numBytes := 0
	for _, gate := range gates {
		if gate[3][0] == 'z' {
			n := parseInt(gate[3][1:])
			numBytes = max(numBytes, n)
		}
	}

	// This is slow and expensive, but we only expect 4 corrections, hence whatever.
	CONN := make(map[string]map[string]map[string]string)
	for _, gate := range gates {
		x, op, y, z := gate[0], gate[1], gate[2], gate[3]
		for _, s := range []string{x, y} {
			if _, ok := CONN[s]; !ok {
				CONN[s] = make(map[string]map[string]string)
			}
			if _, ok := CONN[s][op]; !ok {
				CONN[s][op] = make(map[string]string)
			}
		}
		dest := z
		if alt, ok := subs[z]; ok {
			dest = alt
		}
		CONN[x][op][y] = dest
		CONN[y][op][x] = dest
	}

	carry := CONN["x00"][AND]["y00"]
	assert(CONN["x00"][XOR]["y00"] == "z00", "")

	for b := 1; b < numBytes; b++ {
		x := fmt.Sprintf("x%02d", b)
		y := fmt.Sprintf("y%02d", b)
		z := fmt.Sprintf("z%02d", b)
		xXORy := CONN[x][XOR][y]
		xANDy := CONN[x][AND][y]

		// invariant: x XOR y has XOR
		// x AND y does not have XOR
		if _, ok := CONN[xXORy][XOR]; !ok {
			if _, ok := CONN[xANDy][XOR]; ok {
				debugf("%s swapped with %s", xXORy, xANDy)
				subs[xXORy] = xANDy
				subs[xANDy] = xXORy
				fixGraph(gates, subs)
				return
			} else {
				panic("not implemented")
			}
		}

		// invariant: x XOR y XOR carry -> z
		wantz := CONN[xXORy][XOR][carry]
		if wantz != z {
			subs[wantz] = z
			subs[z] = wantz
			debugf("want %s got %s", z, wantz)
			fixGraph(gates, subs)
			return
		}

		carry = CONN[CONN[xXORy][AND][carry]][OR][xANDy]
	}

	return
}

func dotPlot(gates [][4]string, swaps map[string]string) []byte {
	var b bytes.Buffer
	b.WriteString("digraph G {\n")

	for _, gate := range gates {
		b.WriteString(fmt.Sprintf("  %s_%s_%s[label=%s]\n", gate[0], gate[1], gate[2], gate[1]))
	}

	b.WriteByte('\n')

	for _, gate := range gates {
		b.WriteString(fmt.Sprintf("  %s -> %s_%s_%s\n", gate[0], gate[0], gate[1], gate[2]))
		b.WriteString(fmt.Sprintf("  %s -> %s_%s_%s\n", gate[2], gate[0], gate[1], gate[2]))

		dest := gate[3]
		if swap, ok := swaps[dest]; ok {
			dest = swap
		}

		b.WriteString(fmt.Sprintf("  %s_%s_%s -> %s\n", gate[0], gate[1], gate[2], dest))
	}

	b.WriteString("}")

	return b.Bytes()
}

func main() {
	lines := input()

	init := make(map[string]bool)
	gates := make([][4]string, 0, 1)

	ix := 0
	for ix < len(lines) {
		wire, val := parseWire(lines[ix])
		init[wire] = (val == 1)
		ix++
		if lines[ix] == "" {
			ix++
			break
		}
	}

	for ix < len(lines) {
		in1, typ, in2, out := parseGate(lines[ix])
		gates = append(gates, [4]string{in1, typ, in2, out})
		ix++
	}

	state := compute(init, gates)
	_, _, res1 := getXYZ(state)
	printf("res1: %d", res1)

	subs := make(map[string]string)
	fixGraph(gates, subs)

	swaps := make([]string, 0, len(subs))
	for k := range subs {
		swaps = append(swaps, k)
	}
	sort.Strings(swaps)
	res2 := strings.Join(swaps, ",")
	printf("res2: %s", res2)
}
