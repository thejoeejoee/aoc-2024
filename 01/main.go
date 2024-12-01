package main

import (
	"aoc-2024/internal"
	"fmt"
	"github.com/samber/lo"
	"sort"
	"strings"
)

import "golang.org/x/exp/constraints"

const DemoInput = `
3   4
4   3
2   5
1   3
3   9
3   3
`

var Input string

func init() {
	Input = internal.Download(2024, 1)
}

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	//lines := strings.Split(strings.TrimSpace(DemoInput), "\n")
	lines := strings.Split(strings.TrimSpace(Input), "\n")

	left := make([]int, 0, len(lines))
	right := make([]int, 0, len(lines))

	for _, line := range lines {
		var l, r int
		_ = lo.Must1(fmt.Sscanf(line, "%d %d", &l, &r))
		left = append(left, l)
		right = append(right, r)
	}

	sort.Ints(left)
	sort.Ints(right)

	r := lo.SumBy(lo.Zip2(left, right), func(i lo.Tuple2[int, int]) int {
		return Abs(i.A - i.B)
	})

	fmt.Println(r)

	r = lo.SumBy(left, func(item int) int {
		return item * lo.Count(right, item)
	})

	fmt.Println(r)
}
