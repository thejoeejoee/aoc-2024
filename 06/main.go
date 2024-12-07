package main

import (
	"aoc-2024/internal"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"log/slog"
	"maps"
	"strings"

	_ "aoc-2024/internal"
)

const DemoInput = `
....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...
`

const Loop = `
.#.
#.#
.^.
...
`
const Loop1 = `
...
#^#
.#.
`
const Loop2 = `
#.
^.
`

var Input string

func init() {
	Input = internal.Download(2024, 6)
	//Input = DemoInput
	//Input = Loop1
	//Input = Loop2
	//Input = Loop
}

type coord lo.Tuple2[int, int]

type position struct {
	coord
}

type direction struct {
	coord
}

func (c position) Add(o direction) position {
	return position{coord: coord{
		A: c.A + o.A,
		B: c.B + o.B,
	}}
}

func (c position) In(o position) bool {
	return c.A >= 0 && c.A <= o.A && c.B >= 0 && c.B <= o.B
}

func (d direction) Direction() Flow {
	if d.A != 0 {
		return LeftRight
	}
	return TopDown
}

func (d direction) Turn() direction {
	return directions[(lo.IndexOf(directions, d)+1)%len(directions)]
}

type Flow uint8

const (
	TopDown   Flow = 0b01
	LeftRight Flow = 0b10
)

var Top = direction{coord{A: 0, B: -1}}
var Right = direction{coord{A: 1, B: 0}}
var Bottom = direction{coord{A: 0, B: 1}}
var Left = direction{coord{A: -1, B: 0}}

var directions = []direction{Top, Right, Bottom, Left}

type obstacle = lo.Entry[position, rune]

type State struct {
	dir direction
	pos position
}

//goland:noinspection ALL
func main() {
	state := State{}

	lines := strings.Split(strings.TrimSpace(Input), "\n")

	obstacles := lo.FromEntries(lo.FlatMap(
		lines,
		func(line string, y int) []obstacle {
			return lo.FilterMap([]rune(line), func(ch rune, x int) (obstacle, bool) {
				return obstacle{Key: position{coord{A: x, B: y}}, Value: ch}, ch != '.'
			})
		},
	))

	for i, ch := range "^>v<" {
		pos, ok := lo.Find(lo.Keys(obstacles), func(p position) bool {
			return obstacles[p] == ch
		})
		if ok {
			state.dir = directions[i]
			state.pos = pos

			delete(obstacles, pos)
			break
		}
	}

	height := len(lines)
	width := len(lines[0])
	corner := position{coord{A: width - 1, B: height - 1}}

	lo.Must0(lo.IsNotEmpty(state.pos))
	lo.Must0(lo.IsNotEmpty(state.dir))

	guardInitial := state.pos

	seenPositions := map[position]struct{}{}
	_ = run(state, obstacles, corner, func(s State) {
		//slog.Info("discovered new", "s", s)
		seenPositions[s.pos] = struct{}{}
	})
	fmt.Println(len(seenPositions))

	loopObs := map[position]struct{}{}
	clear(seenPositions)

	_ = run(state, obstacles, corner, func(s State) {
		seenPositions[s.pos] = struct{}{}

		o := s.pos.Add(s.dir)

		switch {
		case o == guardInitial:
			// noop, cannot place on initial
			slog.Info("no obs to guard initial")
			return
		case !o.In(corner):
			// noop, cannot place outside map
			slog.Info("no obs outside map")
			return
		case lo.HasKey(seenPositions, o):
			slog.Info("no obs to known position")
			return
		}

		wNew := maps.Clone(obstacles)
		wNew[o] = 'O'

		err := run(s, wNew, corner, nil)
		if errors.Is(err, ErrInLoop) {
			loopObs[o] = struct{}{}
		}
	})

	fmt.Println(len(loopObs))
}

var (
	ErrInLoop   = errors.New("in loop")
	ErrBoundary = errors.New("boundary")
)

func run(initial State, obs map[position]rune, corner position, discover func(State)) error {
	if discover == nil {
		discover = func(State) {}
	}

	state := initial
	known := map[State]struct{}{}

	record := func(state State) {
		discover(state)
		known[state] = struct{}{}
	}

	for {
		// known state? we're in loop
		if _, loop := known[state]; loop {
			return ErrInLoop
		}

		record(state)

		next := state.pos.Add(state.dir)

		// next one is blocked?
		if _, blocked := obs[next]; blocked {
			state.dir = state.dir.Turn()
			continue
		}

		if !next.In(corner) {
			// next would be outside, break
			return ErrBoundary
		}

		state.pos = next
	}
}

//func debug(o io.Writer, obstacles map[position]rune, s State, seen map[position]Flow, end position, potential *position) {
//	lo.Must(fmt.Fprintln(o, "================="))
//	for y := range end.B + 1 {
//		for x := range end.A + 1 {
//			c := position{coord{A: x, B: y}}
//
//			obs, isObs := obstacles[c]
//			d, wasHere := seen[c]
//
//			switch {
//			case potential != nil && *potential == c:
//				lo.Must(fmt.Fprint(o, "O"))
//			case isObs:
//				lo.Must(fmt.Fprint(o, string(obs)))
//			case c == s.pos:
//				lo.Must(fmt.Fprint(o, string("^>v<"[lo.IndexOf(directions, s.dir)])))
//			case wasHere:
//				lo.Must(fmt.Fprint(o, string(" |-+"[d])))
//			default:
//				lo.Must(fmt.Fprint(o, "."))
//			}
//		}
//		fmt.Print("\n")
//	}
//}
