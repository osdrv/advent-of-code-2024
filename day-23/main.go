package main

import (
	"math/rand"
	"sort"
	"strings"
)

func randPuX(P, X map[string]struct{}) string {
	var PuX map[string]struct{}
	if len(P) == 0 {
		PuX = X
	} else if len(X) == 0 {
		PuX = P
	} else {
		if rand.Intn(1) == 0 {
			PuX = P
		} else {
			PuX = X
		}
	}
	for u := range PuX {
		return u
	}
	panic("wtf")
}

func copyMap[K comparable, V any](m map[K]V) map[K]V {
	cp := make(map[K]V)
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// Implementation of https://en.wikipedia.org/wiki/Bron–Kerbosch_algorithm
func BK(G map[string]map[string]bool, R, P, X map[string]struct{}) map[string]struct{} {
	if len(P) == 0 && len(X) == 0 {
		return copyMap(R)
	}
	var maxclique map[string]struct{}
	pivot := randPuX(P, X)
	for v := range P {
		if G[pivot][v] {
			continue
		}
		R[v] = struct{}{}
		P_ := make(map[string]struct{})
		X_ := make(map[string]struct{})
		for nb := range G[v] {
			if _, ok := P[nb]; ok {
				P_[nb] = struct{}{}
			}
			if _, ok := X[nb]; ok {
				X_[nb] = struct{}{}
			}
		}
		newclique := BK(G, R, P_, X_)
		if len(newclique) > len(maxclique) {
			maxclique = newclique
		}
		delete(R, v)
		delete(P, v)
		X[v] = struct{}{}
	}
	return maxclique
}

func parseConn(s string) (string, string) {
	ss := strings.SplitN(s, "-", 2)
	return ss[0], ss[1]
}

func compute3Sets(conns map[string]map[string]bool) [][]string {
	res := make([][]string, 0, 1)
	uniq := make(map[string]bool)
	for a := range conns {
		if len(conns[a]) < 3 {
			continue
		}
		for b := range conns[a] {
			for c := range conns[a] {
				if !conns[b][c] {
					continue
				}
				if a == b || b == c || a == c {
					continue
				}
				k := []string{a, b, c}
				sort.Strings(k)
				kk := strings.Join(k, "")
				if _, ok := uniq[kk]; !ok {
					res = append(res, k)
					uniq[kk] = true
				}
			}
		}
	}
	return res
}

func main() {
	lines := input()

	G := make(map[string]map[string]bool)
	VX := make(map[string]struct{})
	edges := make([][2]string, 0, len(lines))

	for _, line := range lines {
		a, b := parseConn(line)
		edges = append(edges, [2]string{a, b})
		if _, ok := G[a]; !ok {
			G[a] = make(map[string]bool)
		}
		G[a][b] = true
		if _, ok := G[b]; !ok {
			G[b] = make(map[string]bool)
		}
		G[b][a] = true
		VX[a] = struct{}{}
		VX[b] = struct{}{}
	}

	s3 := compute3Sets(G)
	res1 := 0
NEXTSET:
	for _, s := range s3 {
		for _, ss := range s {
			if strings.HasPrefix(ss, "t") {
				res1++
				continue NEXTSET
			}
		}
	}

	printf("res1: %d", res1)

	clique := BK(G, map[string]struct{}{}, VX, map[string]struct{}{})
	res := make([]string, 0, len(clique))
	for cl := range clique {
		res = append(res, cl)
	}
	sort.Strings(res)
	printf("res2: %s", strings.Join(res, ","))
}
