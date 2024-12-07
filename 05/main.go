package main

import (
	"aoc-2024/internal"
	"fmt"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"math"
	"strconv"
	"strings"
)

const DemoInput = `
47|53
97|13
97|61
97|47
75|29
61|13
75|53
29|13
97|29
53|29
61|53
97|53
61|29
47|13
75|47
97|75
47|61
75|61
47|29
75|13
53|13

75,47,61,53,29
97,61,53,29,13
75,29,13
75,97,47,61,53
61,13,29
97,13,75,29,47
`

var Input string

func init() {
	Input = internal.Download(2024, 5)
	//Input = DemoInput
}

type Rule = lo.Tuple2[int, int]

func main() {
	in := strings.Split(strings.TrimSpace(Input), "\n\n")

	rules := lo.Map(strings.Split(in[0], "\n"), func(item string, _ int) Rule {
		var r Rule
		lo.Must0(2 == lo.Must(fmt.Sscanf(item, "%d|%d", &r.A, &r.B)))
		return r
	})

	c := lo.Sum(parallel.Map(strings.Split(in[1], "\n"), func(line string, _ int) int {
		ints := lo.Map(strings.Split(line, ","), func(n string, _ int) int {
			return lo.Must(strconv.Atoi(n))
		})
		failed := findFailed(ints, rules)
		if len(failed) != 0 {
			return 0
		}

		return ints[len(ints)/2]
	}))

	fmt.Println(c)

	//var swaps int64

	c = lo.Sum(parallel.Map(strings.Split(in[1], "\n"), func(line string, _ int) int {
		ints := lo.Map(strings.Split(line, ","), func(n string, _ int) int {
			return lo.Must(strconv.Atoi(n))
		})

		failed := findFailed(ints, rules)

		if len(failed) == 0 {
			// all fine, we're not computing
			return 0
		}

		for len(failed) > 0 {
			rule := failed[0]
			failed = lo.Slice(failed, 1, math.MaxInt)

			l := lo.IndexOf(ints, rule.A)
			r := lo.IndexOf(ints, rule.B)

			lo.Must0(l >= r, "precondition failed")

			n := ints[l]
			ints = lo.DropByIndex(ints, l)
			ints = lo.Splice(ints, r, n)

			//atomic.AddInt64(&swaps, 1)

			failed = findFailed(ints, rules)
		}

		return ints[len(ints)/2]
	}))

	//slog.Info("swaps", "n", swaps)

	fmt.Println(c)
}

func findFailed(ints []int, rules []Rule) []Rule {
	return lo.FilterMap(rules, func(rule Rule, _ int) (Rule, bool) {
		l := lo.IndexOf(ints, rule.A)
		r := lo.IndexOf(ints, rule.B)

		if l == -1 || r == -1 {
			return Rule{}, false
		}
		return rule, !(r > l)
	})
}
