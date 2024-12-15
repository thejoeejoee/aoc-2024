package main

import (
	"aoc-2024/internal"
	_ "aoc-2024/internal"
	"aoc-2024/internal/taoc"
	"fmt"
	"github.com/sa-/slicefunk"
	"github.com/samber/lo"
	"log/slog"
	"maps"
	"strings"
)

const DemoInput = `
########
#..O.O.#
##@.O..#
#...O..#
#.#.O..#
#...O..#
#......#
########

<^^>>>vv<v>>v<<
`

const DemoInput1 = `
##########
#..O..O.O#
#......O.#
#.OO..O.O#
#..O@..O.#
#O#..O...#
#O..O..O.#
#.OO.O.OO#
#....O...#
##########

<vv>^<v^>v>^vv^v>v<>v^v<v<^vv<<<^><<><>>v<vvv<>^v^>^<<<><<v<<<v^vv^v>^
vvv<<^>^v^^><<>>><>^<<><^vv^^<>vvv<>><^^v>^>vv<>v<<<<v<^v>^<^^>>>^<v<v
><>vv>v^v^<>><>>>><^^>vv>v<^^^>>v^v^<^^>v^^>v^<^v>v<>>v^v^<v>v^^<^^vv<
<<v<^>>^^^^>>>v^<>vvv^><v<<<>^^^vv^<vvv>^>v<^^^^v<>^>vvvv><>>v^<<^^^^^
^><^><>>><>^^<<^^v>>><^<v>^<vv>>v>>>^v><>^v><<<<v>>v<v<v>vvv>^<><<>^><
^>><>^v<><^vvv<^^<><v<<<<<><^v<<<><<<^^<v<^^^><^>>^<v^><<<^>>^v<v^v<v^
>^>>^v>vv>^<<^v<>><<><<v<<v><>v<^vv<<<>^^v^>^^>>><<^v>>v^v><^^>>^<>vv^
<><^^>^^^<><vvvvv^v<v<<>^v<v>v<<^><<><<><<<^^<<<^<<>><<><^^^>^^<>^>v<>
^^>vv<^v^v<vv>^<><v<^v>^^^>>>^^vvv^>vvv<>>>^<^>>>>>^<<^v>^vvv<>^<><<v>
v^^>>><<^^<>>^v^<v^vv<>v^<<>^<^v^v><^<<<><<^<v><v<>vv>>v><v^<vv<>v^<<^
`

const DemoInput2 = `
#######
#...#.#
#.....#
#..OO@#
#..O..#
#.....#
#######

<vv<<^^<<^^
`

var Input string

func init() {
	//Input = DemoInput
	//Input = DemoInput1
	//Input = DemoInput2
	//Input = DemoInputSingle
	Input = internal.Download(2024, 15)
}

type Space struct {
	robot taoc.Position
	walls map[taoc.Position]struct{}
	goods map[taoc.Position]struct{}
	end   taoc.Position
}

type Space2 struct {
	robot taoc.Position
	walls map[taoc.Position]struct{}
	goods map[taoc.Position]rune
	end   taoc.Position
}

func (s Space2) expandToBox(p taoc.Position) []taoc.Position {
	switch s.goods[p] {
	case '[':
		return []taoc.Position{
			p,
			p.Add(taoc.Right),
		}
	case ']':
		return []taoc.Position{
			p.Add(taoc.Left),
			p,
		}
	}
	return []taoc.Position{}
}

func (s Space2) score() int {
	return lo.Sum(lo.Map(lo.Entries(s.goods), func(item lo.Entry[taoc.Position, rune], _ int) int {
		if item.Value == ']' {
			return 0
		}
		return item.Key.A + (item.Key.B * 100)
	}))
}

func (s Space) score() int {
	return lo.Sum(lo.Map(lo.Keys(s.goods), func(item taoc.Position, _ int) int {
		return item.A + (item.B * 100)
	}))
}

func main() {
	//main1()
	main2()
}

func main2() {
	lines := slicefunk.Map(strings.Split(strings.TrimSpace(Input), "\n"), strings.TrimSpace)

	divIdx := lo.IndexOf(lines, "")

	s := Space2{
		walls: map[taoc.Position]struct{}{},
		goods: make(map[taoc.Position]rune),
	}

	lo.ForEach(lines[:divIdx], func(line string, y int) {
		lo.ForEach([]rune(line), func(r rune, x int) {
			p := taoc.NewPos(2*x, y)
			switch r {
			case '#':
				s.walls[p] = struct{}{}
				s.walls[p.Add(taoc.Right)] = struct{}{}
			case 'O':
				s.goods[p] = '['
				s.goods[p.Add(taoc.Right)] = ']'
			case '@':
				s.robot = p
			}
		})
	})

	steps := lo.Map([]rune(strings.Replace(strings.Join(lines[divIdx:], ""), " ", "", -1)), func(r rune, _ int) taoc.Vector {
		return taoc.Vec4FromRune(r)
	})

	s.end.A = len(lines[0])*2 - 1
	s.end.B = divIdx - 1

	debug2(s)
	for _, v := range steps {
		s = step2(s, v)
		//fmt.Printf("after step %s\n", v)
		//debug2(s)
	}

	score := s.score()
	fmt.Println(score)

	//debug(s)
}

//goland:noinspection t
func step2(old Space2, v taoc.Vector) Space2 {
	n := Space2{walls: old.walls, goods: maps.Clone(old.goods), end: old.end, robot: old.robot}

	next := old.robot.Add(v)

	if _, isw := old.walls[next]; isw {
		// noop
		return n
	}
	if _, isg := old.goods[next]; !isg {
		// step into free
		n.robot = next
		return n
	}

	// next one is blocked by good
	if v == taoc.Right || v == taoc.Left {
		// happy path, no stacking

		toPush := []taoc.Position{next}
		blocked := false

		for {
			next = next.Add(v)
			if _, isw := old.walls[next]; isw {
				// is wall
				blocked = true
				break
			}
			if _, isg := old.goods[next]; isg {
				// is good
				toPush = append(toPush, next)
				continue
			}
			// empty
			break
		}
		if blocked {
			return n
		}

		for _, what := range lo.Reverse(toPush) {
			to := what.Add(v)
			n.goods[to] = old.goods[what]
			delete(n.goods, what)
		}
		n.robot = old.robot.Add(v)

		return n
	}

	// only actual boxes, nothing else
	toPush := [][]taoc.Position{old.expandToBox(next)}
	blocked := false

	for {
		allNext := lo.Map(lo.Must(lo.Last(toPush)), func(what taoc.Position, _ int) taoc.Position {
			return what.Add(v)
		})
		allNextBoxes := lo.FlatMap(lo.Must(lo.Last(toPush)), func(what taoc.Position, _ int) []taoc.Position {
			return n.expandToBox(what.Add(v))
		})

		if lo.SomeBy(allNext, lo.Partial(lo.HasKey, old.walls)) {
			// ....#..
			// .[][]..
			// ..[]...
			// ..@....
			blocked = true
			break
		}

		if lo.SomeBy(allNext, lo.Partial(lo.HasKey, old.goods)) {
			// something in next layer
			toPush = append(toPush, allNextBoxes)
			continue
		}

		// nothing in the next layer
		break
	}
	if blocked {
		return n
	}

	allToPush := lo.Reverse(lo.Flatten(toPush))
	slog.Info("multipush", "c", len(allToPush))
	for _, what := range allToPush {
		to := what.Add(v)
		n.goods[to] = old.goods[what]
		delete(n.goods, what)
	}
	n.robot = old.robot.Add(v)

	return n
}

func main1() {
	lines := slicefunk.Map(strings.Split(strings.TrimSpace(Input), "\n"), strings.TrimSpace)

	divIdx := lo.IndexOf(lines, "")

	s := Space{
		walls: map[taoc.Position]struct{}{},
		goods: make(map[taoc.Position]struct{}),
	}

	lo.ForEach(lines[:divIdx], func(line string, y int) {
		lo.ForEach([]rune(line), func(r rune, x int) {
			p := taoc.NewPos(x, y)
			switch r {
			case '#':
				s.walls[p] = struct{}{}
			case 'O':
				s.goods[p] = struct{}{}
			case '@':
				s.robot = p
			}
		})
	})

	steps := lo.Map([]rune(strings.Replace(strings.Join(lines[divIdx:], ""), " ", "", -1)), func(r rune, _ int) taoc.Vector {
		return taoc.Vec4FromRune(r)
	})

	s.end.A = len(lines[0]) - 1
	s.end.B = divIdx - 1

	debug(s)
	for _, v := range steps {
		s = step(s, v)
		//fmt.Printf("after step %s\n", v)
		debug(s)
	}

	var score int = s.score()

	fmt.Println(score)

	//debug(s)
}

func step(old Space, v taoc.Vector) Space {
	n := Space{walls: old.walls, goods: maps.Clone(old.goods), end: old.end, robot: old.robot}

	next := old.robot.Add(v)

	if _, isw := old.walls[next]; isw {
		// noop
		return n
	}
	if _, isg := old.goods[next]; !isg {
		// step into free
		n.robot = next
		return n
	}

	toOccupy := taoc.Position{}
	blocked := false

	for {
		next = next.Add(v)
		if _, isw := old.walls[next]; isw {
			// is wall
			blocked = true
			break
		}
		if _, isg := old.goods[next]; isg {
			// is good
			continue
		}

		toOccupy = next
		// empty
		break
	}
	if blocked {
		return n
	}

	n.robot = old.robot.Add(v)
	delete(n.goods, n.robot)
	n.goods[toOccupy] = struct{}{}

	return n
}

func debug2(s Space2) {
	for y := range s.end.B + 1 {
		for x := range s.end.A + 1 {
			p := taoc.NewPos(x, y)
			g, isg := s.goods[p]
			_, isw := s.walls[p]
			if isg {
				fmt.Print(string(g))
			} else if isw {
				fmt.Print("#")
			} else if p == s.robot {
				fmt.Print("@")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func debug(s Space) {
	for y := range s.end.B + 1 {
		for x := range s.end.A + 1 {
			p := taoc.NewPos(x, y)
			_, isg := s.goods[p]
			_, isw := s.walls[p]
			if isg {
				fmt.Print("O")
			} else if isw {
				fmt.Print("#")
			} else if p == s.robot {
				fmt.Print("@")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
