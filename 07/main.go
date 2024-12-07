package main

import (
	"aoc-2024/internal"
	_ "aoc-2024/internal"
	"fmt"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"strconv"
	"strings"
)

const DemoInput = `
190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20
`

var Input string

func init() {
	Input = internal.Download(2024, 7)
	//Input = DemoInput
}

type Equation lo.Tuple2[int64, []int64]

type opFunc func(a, b int64) (int64, bool)

var ops []opFunc = []opFunc{
	func(a, b int64) (int64, bool) {
		return a - b, true
	},
	func(a, b int64) (int64, bool) {
		if a%b != 0 {
			return 0, false
		}

		return a / b, true
	},
}

var stringOp = func(a, b int64) (int64, bool) {
	as := strconv.FormatInt(a, 10)
	bs := strconv.FormatInt(b, 10)

	possible := strings.HasSuffix(as, bs)

	if !possible || len(as) <= len(bs) || as[0] == '-' {
		//slog.Info("not possible")
		return 0, false
	}

	return lo.Must(strconv.ParseInt(as[:len(as)-len(bs)], 10, 64)), true
}

func (e Equation) resolve() (bool, []Equation) {
	if len(e.B) == 2 {
		//	we're here to solve
		return lo.SomeBy(ops, func(op opFunc) bool {
			r, ok := op(e.A, e.B[0])

			return ok && r == e.B[1]
		}), nil
	}

	return false, lo.FilterMap(ops, func(op opFunc, _ int) (Equation, bool) {
		r, ok := op(e.A, e.B[0])

		return Equation{
			A: r,
			B: e.B[1:],
		}, ok

	})
}

func main() {
	eqs := lo.Map(strings.Split(strings.TrimSpace(Input), "\n"), func(l string, _ int) Equation {
		ns := lo.Map(strings.Fields(strings.ReplaceAll(l, ":", " ")), func(n string, index int) int64 {
			return lo.Must(strconv.ParseInt(n, 10, 64))
		})

		return Equation{
			A: ns[0],
			// to respect left-to-right w/o getting headache
			B: lo.Reverse(ns[1:]),
		}
	})

	sum := solve(eqs)
	fmt.Println(sum)

	ops = append(ops, stringOp)

	sum = solve(eqs)
	fmt.Println(sum)
}

func solve(eqs []Equation) int64 {
	solvable := parallel.Map(eqs, func(root Equation, _ int) *Equation {
		q := []Equation{root}

		for {
			if len(q) == 0 {
				return nil
			}
			eq := q[0]
			q = q[1:]

			resolved, rest := eq.resolve()
			if resolved {
				return &root
			}

			q = append(q, rest...)
		}
	})

	return lo.SumBy(solvable, func(eq *Equation) int64 {
		if eq == nil {
			return 0
		}
		return eq.A
	})
}
