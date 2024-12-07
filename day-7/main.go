package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Calibr struct {
	Val  uint64
	Nums []uint64
}

func (c *Calibr) String() string {
	return fmt.Sprintf("C{Val: %d, Nums: %+v}", c.Val, c.Nums)
}

type qitem struct {
	ix  int
	res uint64
}

type Ops uint8

const (
	OpAdd = 1 << iota
	OpMul
	OpConcat
)

func (c *Calibr) IsTrue(ops Ops) bool {
	q := make([]qitem, 0, 1)
	q = append(q, qitem{1, c.Nums[0]})
	var head qitem
	for len(q) > 0 {
		head, q = q[0], q[1:]
		if head.ix == len(c.Nums) {
			if head.res == c.Val {
				return true
			}
			continue
		}
		// all nums are unsigned, we can not reach the result
		if head.res > c.Val {
			continue
		}
		if ops&OpAdd != 0 {
			nres := head.res + c.Nums[head.ix]
			q = append(q, qitem{head.ix + 1, nres})
		}
		if ops&OpMul != 0 {
			nres := head.res * c.Nums[head.ix]
			q = append(q, qitem{head.ix + 1, nres})
		}
		if ops&OpConcat != 0 {
			nds := numDigs(c.Nums[head.ix])
			nres := head.res
			for i := 1; i <= nds; i++ {
				nres *= 10
			}
			nres += c.Nums[head.ix]
			q = append(q, qitem{head.ix + 1, nres})
		}
	}

	return false
}

func ParseCalibr(s string) (*Calibr, error) {
	ss := strings.SplitN(s, ": ", 2)
	if len(ss) != 2 {
		return nil, fmt.Errorf("invalid calibr: %s", s)
	}
	val, err := strconv.ParseUint(ss[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse calibr value: %s", err)
	}
	nums := make([]uint64, 0, 1)
	for _, s := range strings.Split(ss[1], " ") {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse calibr number: %s", err)
		}
		nums = append(nums, n)
	}
	return &Calibr{Val: val, Nums: nums}, nil
}

func main() {
	lines := input()
	cs := make([]*Calibr, 0, len(lines))
	for _, line := range lines {
		c, err := ParseCalibr(line)
		noerr(err)
		cs = append(cs, c)
		debugf("c: %+v", cs)
	}

	sum1 := uint64(0)
	for _, c := range cs {
		if c.IsTrue(OpAdd | OpMul) {
			sum1 += c.Val
		}
	}
	printf("sum1: %d", sum1)

	sum2 := uint64(0)
	for _, c := range cs {
		if c.IsTrue(OpAdd | OpMul | OpConcat) {
			sum2 += c.Val
		}
	}
	printf("sum2: %d", sum2)
}
