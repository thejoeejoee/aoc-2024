package main

import (
	"aoc-2024/internal"
	_ "aoc-2024/internal"
	"aoc-2024/internal/taoc"
	"fmt"
	"github.com/sa-/slicefunk"
	"github.com/samber/lo"
	//"iter"
	"log/slog"
	"math"
	"slices"
	"strings"
)

const DemoInput = `
###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############
`
const DemoInput1 = `
#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################
`

var Input string

func init() {
	//Input = DemoInput1
	//Input = DemoInputSingle
	Input = internal.Download(2024, 16)
}

type WallEntry = lo.Entry[taoc.Position, struct{}]

type orientedPos struct {
	pos taoc.Position
	dir taoc.Vector
}

type bestPath struct {
	path  []taoc.Position
	score int
}

type State struct {
	reindeer   orientedPos
	walls      map[taoc.Position]struct{}
	target     taoc.Position
	best       map[taoc.Position]map[taoc.Vector]int
	bestPaths  map[orientedPos]bestPath
	toDiscover []ToDiscover

	end taoc.Position
}

func (s State) shouldDiscover(opos orientedPos, newScore int) bool {
	if _, isw := s.walls[opos.pos]; isw {
		return false
	}

	if current, seen := s.best[opos.pos][opos.dir]; seen && newScore > current {
		return false
	}

	return true

}

type ToDiscover struct {
	by    []taoc.Position
	dir   taoc.Vector
	score int
}

//goland:noinspection Flavor,t
func main() {
	lines := slicefunk.Map(strings.Split(strings.TrimSpace(Input), "\n"), strings.TrimSpace)

	s := State{}

	s.end = taoc.NewPos(len(lines[0])-1, len(lines)-1)

	s.walls = lo.FromEntries(lo.FlatMap(lines, func(line string, y int) []WallEntry {
		return lo.FilterMap([]rune(line), func(r rune, x int) (WallEntry, bool) {
			p := taoc.NewPos(x, y)

			switch r {
			case '#':
				return WallEntry{Key: taoc.NewPos(x, y)}, true
			case 'S':
				s.reindeer.pos = p
				s.reindeer.dir = taoc.Right
			case 'E':
				s.target = p
			}

			return WallEntry{}, false
		})
	}))

	s.best = map[taoc.Position]map[taoc.Vector]int{s.reindeer.pos: {s.reindeer.dir: 0}}
	s.bestPaths = map[orientedPos]bestPath{}

	s.toDiscover = make([]ToDiscover, 0, 128)
	s.toDiscover = append(s.toDiscover, ToDiscover{[]taoc.Position{s.reindeer.pos}, s.reindeer.dir, 0})

	i := 0
	for len(s.toDiscover) > 0 {
		//fmt.Println("======================")
		td := s.toDiscover[0]
		s.toDiscover = lo.Slice(s.toDiscover, 1, math.MaxInt)
		curr := lo.Must(lo.Last(td.by))
		curropos := orientedPos{pos: curr, dir: td.dir}

		//slog.Info("discovering", "curr", curr, "dir", td.dir.String())
		i++
		if i%100_000 == 0 {
			slog.Info("discovered", "i", i, "len", len(s.toDiscover))
		}
		//slog.Info("current", "i", i, "len", len(s.toDiscover))
		//slog.Info("current", "curr", curr)
		//debug(s, td)

		if _, f := s.best[curropos.pos]; !f {
			s.best[curropos.pos] = map[taoc.Vector]int{}
		}
		if current, has := s.best[curropos.pos][curropos.dir]; !has || td.score < current {
			s.best[curropos.pos][curropos.dir] = td.score
		}

		if bp, has := s.bestPaths[curropos]; !has {
			s.bestPaths[curropos] = bestPath{
				path:  slices.Clone(td.by),
				score: td.score,
			}
		} else {
			if bp.score == td.score {
				bp.path = append(s.bestPaths[curropos].path, slices.Clone(td.by)...)
				s.bestPaths[curropos] = bp
			} else if td.score < bp.score {
				s.bestPaths[curropos] = bestPath{
					path:  slices.Clone(td.by),
					score: td.score,
				}
			}
		}

		if curr == s.target {
			// do not discover more
			continue
		}

		for _, off := range []int{-1, 1} {
			v := lo.Must(lo.Nth(taoc.Directions4, (lo.IndexOf(taoc.Directions4, td.dir)+off)%len(taoc.Directions4)))

			nextAfterTurn := curr.Add(v)
			opos := orientedPos{pos: nextAfterTurn, dir: v}

			if s.shouldDiscover(opos, td.score+1001) {
				//slog.Info("to discover w turn", "p", nextAfterTurn, "dir", v.String())
				by := slices.Clone(td.by)

				s.toDiscover = lo.Splice(s.toDiscover, 0, ToDiscover{
					by:    append(by, nextAfterTurn),
					dir:   v,
					score: td.score + 1000 + 1,
				})
			}
		}

		nextStraight := curr.Add(td.dir)
		movedopos := orientedPos{pos: nextStraight, dir: td.dir}
		if s.shouldDiscover(movedopos, td.score+1) {
			//slog.Info("to discover straight", "p", nextStraight, "dir", td.dir.String())
			by := slices.Clone(td.by)

			s.toDiscover = lo.Splice(s.toDiscover, 0, ToDiscover{
				by:    append(by, nextStraight),
				dir:   td.dir,
				score: td.score + 1,
			})
		}

	}

	slog.Info("discovered", "best", s.best)
	slog.Info("target", "targets", s.best[s.target])

	for _, dir := range taoc.Directions4 {
		opos := orientedPos{pos: s.target, dir: dir}
		slog.Info("best positions", "bp", len(lo.Uniq(s.bestPaths[opos].path)))
	}

}

func debug(s State, td ToDiscover) {
	last := lo.Must(lo.Last(td.by))
	for y := range s.end.B + 1 {
		for x := range s.end.A + 1 {
			p := taoc.NewPos(x, y)
			_, isw := s.walls[p]
			//_, iss := s.best[p]

			if isw {
				fmt.Print("#")
			} else if p == s.target && last == p {
				fmt.Print("!")
			} else if p == last {
				fmt.Print(td.dir.String())
			} else if i := lo.IndexOf(td.by, p); i != -1 {
				fmt.Print("@")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
