package main

import (
	"aoc-2024/internal"
	_ "aoc-2024/internal"
	"aoc-2024/internal/taoc"
	"fmt"
	"github.com/alex-ant/gomath/gaussian-elimination"
	"github.com/alex-ant/gomath/rational"
	"github.com/sa-/slicefunk"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"log/slog"
	"strings"
)

const DemoInput = `
Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279
`

var Input string

func init() {
	//Input = DemoInput
	Input = internal.Download(2024, 13)
}

type Machine struct {
	buttons map[int]taoc.Vector
	target  taoc.Position
}

const (
	ButtonA = 3
	ButtonB = 1
)

func main() {
	sections := slicefunk.Map(strings.Split(strings.TrimSpace(Input), "\n\n"), strings.TrimSpace)

	machines := lo.Map(sections, func(section string, _ int) Machine {

		a := taoc.Vector{}
		b := taoc.Vector{}
		t := taoc.Position{}

		// Button A: X+69, Y+23
		// Button B: X+27, Y+71
		// Prize: X=18641, Y=10279
		lo.Must(fmt.Sscanf(
			section,
			"Button A: X+%d, Y+%d\n"+
				"Button B: X+%d, Y+%d\n"+
				"Prize: X=%d, Y=%d\n",
			&a.A, &a.B,
			&b.A, &b.B,
			&t.A, &t.B,
		))

		m := Machine{
			target:  t,
			buttons: map[int]taoc.Vector{ButtonA: a, ButtonB: b},
		}

		return m
	})

	slog.Info("original")
	c := lo.Sum(parallel.Map(machines, func(m Machine, _ int) int {
		p := solveForMachine(m)
		slog.Info("solvable?", "p", p != nil)
		return lo.FromPtrOr(p, 0)
	}))

	fmt.Println(c)

	slog.Info("boosted")
	for i := range machines {
		machines[i].target.A += 10000000000000
		machines[i].target.B += 10000000000000
	}

	c = lo.Sum(parallel.Map(machines, func(m Machine, _ int) int {
		p := solveSystemForMachine(m)
		slog.Info("solvable?", "p", p != nil)
		return lo.FromPtrOr(p, 0)
	}))

	fmt.Println(c)
}

func solveSystemForMachine(m Machine) *int {
	// N/M, count of A/B
	// AA A.a, ...
	// TA, target.A

	// TA known, AA, AB known
	//
	// 0 = N*_AA + M*_BA - _TA
	// 0 = N*_AB + M*_BB - _TB

	// looking for min(3M + 1N)
	nr := func(i int) rational.Rational {
		return rational.New(int64(i), 1)
	}

	aa := m.buttons[ButtonA].A
	ab := m.buttons[ButtonA].B
	ba := m.buttons[ButtonB].A
	bb := m.buttons[ButtonB].B

	equations := [][]rational.Rational{
		{nr(aa), nr(ba), nr(m.target.A)},
		{nr(ab), nr(bb), nr(m.target.B)},
	}

	sols := lo.Must(gaussian.SolveGaussian(equations, false))

	lo.Must0(len(sols) == 2)
	lo.Must0(len(sols[0]) == 1)
	lo.Must0(len(sols[1]) == 1)

	// for each possible solution, calculate the price
	a := sols[0][0]
	b := sols[1][0]

	if !(a.IsNatural() && b.IsNatural()) {
		return nil
	}

	a.Simplify()
	b.Simplify()

	return lo.ToPtr(int(a.GetNumerator())*ButtonA + int(b.GetNumerator())*ButtonB)

}

func solveForMachine(m Machine) *int {
	forA := lo.Must(allSums(m.target.A, m.buttons[ButtonA].A, m.buttons[ButtonB].A))
	forB := lo.Must(allSums(m.target.B, m.buttons[ButtonA].B, m.buttons[ButtonB].B))

	possibles := lo.Intersect(forA, forB)

	if len(possibles) == 0 {
		return nil
	}

	prices := lo.Map(possibles, func(p lo.Tuple2[int, int], _ int) int {
		return p.A*ButtonA + p.B*ButtonB
	})

	//slog.Info("possibles", "m", m.target, "p", possibles)
	return lo.ToPtr(lo.Min(prices))
}

func allSums(target, a, b int) (results []lo.Tuple2[int, int], err error) {
	if target < 0 || a < 0 || b < 0 {
		return nil, fmt.Errorf("negative values not allowed")
	}

	for i := 0; i <= target/a; i++ {
		rest := target - i*a
		if rest%b == 0 {
			results = append(results, lo.Tuple2[int, int]{A: i, B: rest / b})
		}
	}

	return
}
