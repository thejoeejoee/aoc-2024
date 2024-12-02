package main

import (
	"aoc-2024/internal"
	"fmt"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"strconv"
	"strings"
)

const DemoInput = `
7 6 4 2 1
1 2 7 8 9
9 7 6 2 1
1 3 2 4 5
8 6 4 4 1
1 3 6 7 9
`

var Input string = DemoInput

func init() {
	Input = internal.Download(2024, 2)
	//Input = DemoInput
}

func isOk(r []int) bool {
	diffs := lo.ZipBy2(
		lo.Slice(r, 0, len(r)-1),
		lo.Slice(r, 1, len(r)),
		func(a, b int) int {
			return a - b
		},
	)

	allNegative := lo.EveryBy(diffs, func(d int) bool {
		return d < 0
	})
	allPositive := lo.EveryBy(diffs, func(d int) bool {
		return d > 0
	})
	minDiff := lo.Min(diffs)
	maxDiff := lo.Max(diffs)

	return (allNegative && minDiff >= -3) || (allPositive && maxDiff <= 3)

}

func main() {
	reports := lo.Map(strings.Split(strings.TrimSpace(Input), "\n"), func(r string, _ int) []int {
		return lo.Map(strings.Fields(r), func(f string, _ int) int {
			return lo.Must(strconv.Atoi(f))
		})
	})
	ok := lo.CountBy(
		reports,
		isOk,
	)

	fmt.Println(ok)

	ok = lo.CountBy(reports, func(item []int) bool {
		return lo.SomeBy(parallel.Map(item, func(_ int, index int) bool {
			return isOk(lo.DropByIndex(item, index))
		}), func(item bool) bool {
			return item
		})
	})

	fmt.Println(ok)
}
