package main

import (
	"aoc-2024/internal"
	_ "aoc-2024/internal"
	"aoc-2024/internal/taoc"
	"fmt"
	"github.com/sa-/slicefunk"
	"github.com/samber/lo"
	"maps"
	"strconv"
	"strings"
)

const DemoInput = `
p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3
`

const DemoInputSingle = `
p=2,4 v=2,-3
`

var Input string

func init() {
	Input = DemoInput
	//Input = DemoInputSingle
	Input = internal.Download(2024, 14)
}

type Robot struct {
	taoc.Position
	taoc.Vector
}

type Space struct {
	counts  map[taoc.Position]int
	vectors map[taoc.Position][]taoc.Vector
	end     taoc.Position

	robots int
}

func (s Space) clone() Space {
	return Space{
		counts:  make(map[taoc.Position]int, s.robots),
		vectors: make(map[taoc.Position][]taoc.Vector, s.robots),
		end:     s.end,
	}
}

func (s Space) safety() int {
	// 0 1 X 3 4 -> >2
	// X X 3 X X X -> >3
	midw := s.end.A / 2
	midh := s.end.B / 2

	for x := range s.end.A + 1 {
		s.counts[taoc.NewPos(x, midh)] = 0
	}
	for y := range s.end.B + 1 {
		s.counts[taoc.NewPos(midw, y)] = 0
	}

	debug(s)

	qs := map[lo.Tuple2[bool, bool]]int{}

	lo.ForEach(lo.Entries(s.counts), func(p lo.Entry[taoc.Position, int], _ int) {
		qs[lo.Tuple2[bool, bool]{
			A: p.Key.A > midw,
			B: p.Key.B > midh,
		}] += p.Value
	})

	//fmt.Println(qs)

	return lo.Reduce(lo.Values(qs), func(agg int, c int, _ int) int {
		return agg * c
	}, 1)
}

//goland:noinspection Flavor,t
func main() {
	lines := slicefunk.Map(strings.Split(strings.TrimSpace(Input), "\n"), strings.TrimSpace)
	s := Space{robots: len(lines)}.clone()

	lo.ForEach(lines, func(line string, _ int) {
		r := Robot{}

		lo.Must(fmt.Sscanf(
			line,
			"p=%d,%d v=%d,%d",
			&r.Position.A, &r.Position.B, &r.Vector.A, &r.Vector.B,
		))

		s.counts[r.Position] += 1
		s.vectors[r.Position] = append(
			s.vectors[r.Position],
			r.Vector,
		)
	})
	s.end = taoc.Position{Coord: taoc.Coord{
		A: lo.Max(lo.Map(lo.Keys(s.counts), func(p taoc.Position, _ int) int {
			return p.A
		})),
		B: lo.Max(lo.Map(lo.Keys(s.counts), func(p taoc.Position, _ int) int {
			return p.B
		})),
	}}
	s.end = taoc.NewPos(100, 102)

	original := s.clone()
	original.counts = maps.Clone(s.counts)
	original.vectors = maps.Clone(s.vectors)

	//debug(s)

	for i := range 100 {
		s = step(s)
		_ = i
		//debug(s)
	}

	c := s.safety()
	fmt.Println(c)

	s = original
	i := 0
	for {
		i++
		s = step(s)
		fmt.Printf("this step? %d\n", i)

		// no overlaps
		if lo.Max(lo.Values(s.counts)) == 1 {
			break
		}
	}
	debug(s)
}

func step(s Space) Space {
	n := s.clone()
	lo.ForEach(lo.Entries(s.vectors), func(cell lo.Entry[taoc.Position, []taoc.Vector], _ int) {
		lo.ForEach(cell.Value, func(v taoc.Vector, _ int) {
			jumped := cell.Key.Add(v)
			jumped = jumped.WrapAround(s.end)

			n.vectors[jumped] = append(
				n.vectors[jumped],
				v,
			)
			n.counts[jumped] += 1
		})
	})

	return n
}

func debug(s Space) {
	for y := range s.end.B + 1 {
		for x := range s.end.A + 1 {
			c, has := s.counts[taoc.NewPos(x, y)]
			if has {
				fmt.Print(strconv.Itoa(c))
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
