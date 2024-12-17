package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ADV = 0
	BXL = 1
	BST = 2
	JNZ = 3
	BXC = 4
	OUT = 5
	BDV = 6
	CDV = 7
)

type Register struct {
	A, B, C int
}

func combo(R Register, op int) int {
	switch op {
	case 0, 1, 2, 3:
		return op
	case 4:
		return R.A
	case 5:
		return R.B
	case 6:
		return R.C
	case 7:
		panic("But Eric promised!")
	default:
		panic("wtf")
	}
}

func compute(R Register, prog []int) []int {
	out := make([]int, 0, 1)
	ip := 0

	for {
		if ip >= len(prog) {
			debugf("HALT")
			break
		}
		instr := prog[ip]
		op := prog[ip+1]

		switch instr {
		case ADV: // 0
			res := R.A / (1 << combo(R, op))
			R.A = res
			ip += 2
		case BXL: // 1
			res := R.B ^ op
			R.B = res
			ip += 2
		case BST: // 2
			res := combo(R, op) % 8
			R.B = res
			ip += 2
		case JNZ: // 3
			if R.A == 0 {
				ip += 2
				break
			}
			ip = op // instruction pointer is not incremented but substituted
		case BXC: // 4
			res := R.B ^ R.C
			R.B = res
			ip += 2
		case OUT: // 5
			res := combo(R, op) % 8
			out = append(out, res)
			ip += 2
		case BDV: // 6
			res := R.A / (1 << combo(R, op))
			R.B = res
			ip += 2
		case CDV: // 7
			res := R.A / (1 << combo(R, op))
			R.C = res
			ip += 2
		default:
			panic("wtf")
		}
	}

	return out
}

func printOut(nums []int) string {
	var s strings.Builder
	for _, n := range nums {
		if s.Len() > 0 {
			s.WriteByte(',')
		}
		s.WriteString(strconv.Itoa(n))
	}
	return s.String()
}

func main() {
	lines := input()
	var R Register
	fmt.Sscanf(lines[0], "Register A: %d", &R.A)
	fmt.Sscanf(lines[1], "Register B: %d", &R.B)
	fmt.Sscanf(lines[2], "Register C: %d", &R.C)

	prog := parseInts(strings.SplitN(lines[4], ": ", 2)[1])

	out := compute(R, prog)
	printf("res1: %s", printOut(out))

	mknum := func(nums []int) int {
		res := 0
		for i := 0; i < len(nums); i++ {
			res *= 8
			res += nums[i]
		}
		return res
	}

	nums := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0}

DELTA:
	for d := 0; d < len(prog)-2; d++ {
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
			STEP_K:
				for k := 0; k < 8; k++ {
					nums[d] = i
					nums[d+1] = j
					nums[d+2] = k
					num := mknum(nums)
					res := compute(Register{num, 0, 0}, prog)
					if len(res) != len(prog) {
						continue
					}
					for b := len(nums) - 1; b > len(nums)-1-d-3; b-- {
						if prog[b] != res[b] {
							continue STEP_K
						}
					}
					nums[d] = i
					debugf("i=%d j=%d k=%d mknum0: %d", i, j, k, num)
					debugf("res: %+v", res)
					debugf("nums: %+v", nums)
					continue DELTA
				}
			}
		}
	}

	printf("final nums: %+v", nums)
	printf("res2: %d", mknum(nums))
}
