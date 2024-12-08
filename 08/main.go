package main

import (
	"aoc-2024/internal"
	"fmt"
	"github.com/mowshon/iterium"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"io"
	"log/slog"
	"math"
	"strings"

	_ "aoc-2024/internal"
)

const DemoInput = `
..........
..........
..........
....a.....
..........
.....a....
..........
..........
..........
..........
`
const DemoInput2 = `
..........
..........
..........
....a.....
........a.
.....a....
..........
..........
..........
..........
`

const DemoInput3 = `
..........
..........
..........
....a.....
........a.
.....a....
..........
......A...
..........
..........
`
const DemoInput4 = `
............
........0...
.....0......
.......0....
....0.......
......A.....
............
............
........A...
.........A..
............
............
`
const DemoInput5 = `
T.........
...T......
.T........
..........
..........
..........
..........
..........
..........
..........
`

var Input string

func init() {
	Input = DemoInput
	Input = DemoInput2
	Input = DemoInput3
	Input = DemoInput4
	Input = internal.Download(2024, 8)
}

type coord lo.Tuple2[int, int]

type position struct {
	coord
}

type vector struct {
	coord
}

func (c position) Add(o vector) position {
	return position{coord: coord{
		A: c.A + o.A,
		B: c.B + o.B,
	}}
}

func (c position) Diff(o position) vector {
	return vector{coord: coord{
		A: c.A - o.A,
		B: c.B - o.B,
	}}
}

func (c position) In(o position) bool {
	return c.A >= 0 && c.A <= o.A && c.B >= 0 && c.B <= o.B
}

type antenna = lo.Entry[position, rune]

func main() {
	lines := strings.Split(strings.TrimSpace(Input), "\n")

	antennas := lo.FromEntries(lo.FlatMap(
		lines,
		func(line string, y int) []antenna {
			return lo.FilterMap([]rune(line), func(ch rune, x int) (antenna, bool) {
				return antenna{Key: position{coord{A: x, B: y}}, Value: ch}, ch != '.'
			})
		},
	))

	corner := position{coord{A: len(lines[0]) - 1, B: len(lines) - 1}}

	//debug(os.Stderr, antennas, corner)
	//for _, a := range antinodes {
	//	if _, ok := antennas[a]; !ok {
	//		antennas[a] = '#'
	//	}
	//}
	//debug(os.Stderr, antennas, corner)

	antinodes := findAntinodes(antennas, corner, 1)
	fmt.Println(len(antinodes))

	antinodes = findAntinodes(antennas, corner, math.MaxInt)
	fmt.Println(len(antinodes))

	//for _, a := range antinodes {
	//	if _, ok := antennas[a]; !ok {
	//		antennas[a] = '#'
	//	}
	//}
	//debug(os.Stderr, antennas, corner)
}

func findAntinodes(antennas map[position]rune, corner position, limit int) []position {
	groups := lo.GroupBy(lo.Keys(antennas), func(a position) rune {
		return antennas[a]
	})

	antinodes := parallel.Map(lo.Entries(groups), func(item lo.Entry[rune, []position], index int) []position {
		// all possible 2-tuples of antennas
		tuples := iterium.Permutations(item.Value, 2)

		antinodes := lo.Flatten(lo.Map(lo.Must(tuples.Slice()), func(tuple []position, _ int) []position {
			b := tuple[1]
			a := tuple[0]
			v := b.Diff(a)

			found := []position{b}
			for i := 0; i < limit; i++ {
				b = b.Add(v)

				if !b.In(corner) {
					return found
				}

				found = append(found, b)
			}

			slog.Info("new ans", "a", a, "found", found)
			return found
		}))

		return antinodes
	})

	return lo.Uniq(lo.Flatten(antinodes))
}

func debug(o io.Writer, antennas map[position]rune, corner position) {
	for y := range corner.B + 1 {
		for x := range corner.A + 1 {
			c := position{coord{A: x, B: y}}

			a, isA := antennas[c]

			switch {
			case isA:
				lo.Must(fmt.Fprint(o, string(a)))
			default:
				lo.Must(fmt.Fprint(o, "."))
			}
		}
		lo.Must(fmt.Fprint(o, "\n"))
	}
	lo.Must(fmt.Fprint(o, "\n"))
}
