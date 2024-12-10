package main

import (
	_ "aoc-2024/internal"
	"fmt"
	"github.com/alitto/pond/v2"
	"github.com/sa-/slicefunk"
	"github.com/samber/lo"
	"log/slog"
	"strings"
)

const DemoInput = `
10..9..
2...8..
3...7..
4567654
...8..3
...9..2
.....01
`
const DemoInput2 = `
89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732
`
const DemoInput3 = `
012345
123456
234567
345678
4.6789
56789.
`

var Input string

func init() {
	Input = DemoInput3
	//Input = internal.Download(2024, 10)
}

var Top = vector{coord{A: 0, B: -1}}
var Right = vector{coord{A: 1, B: 0}}
var Bottom = vector{coord{A: 0, B: 1}}
var Left = vector{coord{A: -1, B: 0}}

var directions = []vector{Top, Right, Bottom, Left}

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

func (c position) In(end position) bool {
	return c.A >= 0 && c.A <= end.A && c.B >= 0 && c.B <= end.B
}

type Hill = lo.Entry[position, uint8]

type route struct {
	seen []position
}

func (r route) completed() bool {
	return len(r.seen) == 10
}

type state struct {
	end   position
	hills map[position]uint8
	p     pond.ResultPool[[]route]

	todo    []route
	results []route
}

func main() {
	lines := slicefunk.Map(strings.Split(strings.TrimSpace(Input), "\n"), strings.TrimSpace)
	hills := lo.FromEntries(lo.FlatMap(lines, func(line string, y int) []Hill {
		return lo.Map([]rune(line), func(h rune, x int) Hill {
			return Hill{Key: position{coord{A: x, B: y}}, Value: uint8(h - '0')}
		})
	}))

	end := position{coord{A: len(lines[0]), B: len(lines)}}

	//fmt.Println(hills)

	initials := lo.FilterMap(lo.Entries(hills), func(h Hill, _ int) (route, bool) {
		if h.Value != 0 {
			return route{}, false
		}
		return route{seen: []position{h.Key}}, true
	})

	//fmt.Println(initials)

	//pool := pond.NewResultPool[[]route](0)
	//g := pool.NewGroup()

	s := &state{end: end, hills: hills, todo: []route{}}

	s.todo = append(s.todo, initials...)

	for len(s.todo) > 0 {
		r := s.todo[0]
		sliced := s.todo[1:]
		news := discover(s, r)
		slog.Info("discovered", "news", lo.Map(news, func(item route, index int) position {
			return lo.Must(lo.Last(item.seen))
		}))
		s.todo = append(sliced, news...)
	}

	results := s.results

	c := lo.SumBy(initials, func(i route) int {
		all := lo.Filter(results, func(r route, _ int) bool {
			return r.seen[0] == i.seen[0]
		})
		uniqEnds := lo.UniqBy(all, func(r route) position {
			return lo.Must(lo.Last(r.seen))
		})
		for _, r := range all {
			fmt.Println(i, "===", r)
		}

		return len(uniqEnds)
	})

	fmt.Println(c)

	fmt.Println(len(results))
}

func discover(s *state, r route) []route {
	log := slog.With("r", r)

	if r.completed() {
		// finalize if completed
		log.Info("completed")
		s.results = append(s.results, r)
		return nil
	}

	return lo.FilterMap(directions, func(dir vector, _ int) (route, bool) {
		//	try direction and return new route if we want to continue
		current := lo.Must(lo.Last(r.seen))
		next := current.Add(dir)

		log := log.With("dir", dir, "current", current, "next", next)

		if !next.In(s.end) {
			log.Info("nope, outside")
			// outside
			return route{}, false
		}

		if !(s.hills[next] == s.hills[current]+1) {
			log.Info("nope, not uphill")
			return route{}, false
		}

		expanded := route{seen: append([]position{}, r.seen...)}
		expanded.seen = append(expanded.seen, next)

		log.Info("oh yes!", "expanded", expanded)
		return expanded, true
	})
}
