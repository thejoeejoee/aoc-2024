package main

import (
	"aoc-2024/internal"
	"fmt"
	"github.com/samber/lo"
	"strings"
)

const DemoInput1 = `
ZUGAGDSUZXSD
XMASXMASXMAS
MMIUOHAUDASD
AYAADSASDSDS
SERSASDASDDD
`
const DemoInput = `
MMMSXXMASM
MSAMXMSMSA
AMXSXMAAMM
MSAMASMSMX
XMASAMXAMM
XXAMMXXAMA
SMSMSASXSS
SAXAMASAAA
MAMMMXMMMM
MXMXAXMASX
`

const x = `
.M.S......
..A..MSMS.
.M.S.MAA..
..A.ASMSM.
.M.S.M....
..........
S.S.S.S.S.
.A.A.A.A..
M.M.M.M.M.
..........
`

const Verb = "XMAS"
const VerbLen = len(Verb)

var Input string

func init() {
	Input = internal.Download(2024, 4)
	//Input = DemoInput
}

type coord lo.Tuple2[int, int]

func (c coord) Add(o coord) coord {
	return coord{
		A: c.A + o.A,
		B: c.B + o.B,
	}
}

var TopLeft = coord{A: -1, B: -1}
var TopRight = coord{A: 1, B: -1}
var BottomLeft = coord{A: -1, B: 1}
var BottomRight = coord{A: 1, B: 1}

var directions = []coord{
	{0, 1},
	{0, -1},
	BottomRight,
	{1, 0},
	TopRight,
	BottomLeft,
	{-1, 0},
	TopLeft,
}

func main() {
	type boundRune = lo.Entry[coord, rune]

	m := lo.FromEntries(lo.FlatMap(
		strings.Split(strings.TrimSpace(Input), "\n"),
		func(l string, y int) []boundRune {
			return lo.Map([]rune(l), func(r rune, x int) boundRune {
				return boundRune{Key: coord{A: x, B: y}, Value: r}
			})
		}),
	)

	combs := lo.FlatMap(lo.Keys(m), func(base coord, _ int) [][]coord {
		return lo.Map(directions, func(dir coord, _ int) []coord {
			return coords(base, dir, VerbLen)
		})
	})

	c := lo.CountBy(combs, func(ixs []coord) bool {
		return strings.Join(lo.Map(ixs, func(c coord, _ int) string {
			r, ok := m[c]
			if ok {
				return string(r)
			}
			return "."
		}), "") == Verb
	})

	fmt.Println(c)

	as := lo.Filter(lo.Entries(m), func(i boundRune, _ int) bool {
		return i.Value == 'A'
	})

	var pairs = [][2]coord{
		{TopRight, BottomLeft},
		{TopLeft, BottomRight},
		{BottomLeft, TopRight},
		{BottomRight, TopLeft},
	}

	c = lo.CountBy(as, func(root boundRune) bool {
		return lo.CountBy(pairs, func(item [2]coord) bool {
			l := item[0]
			r := item[1]

			return m[root.Key.Add(l)] == 'M' && m[root.Key.Add(r)] == 'S'
		}) == 2 // two matches per the cross, not just one
	})
	fmt.Println(c)
}

func coords(base coord, dir coord, length int) []coord {
	offsets := lo.Times(length-1, func(i int) int {
		return i
	})
	return Accumulate(offsets, func(agg coord, item int, _ int) coord {
		return agg.Add(dir)
	}, base)
}

func Accumulate[T any, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) []R {
	var out []R

	for i := range collection {
		out = append(out, initial)
		initial = accumulator(initial, collection[i], i)
	}

	out = append(out, initial)

	return out
}
