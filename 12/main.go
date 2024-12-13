package main

import (
	"aoc-2024/internal"
	_ "aoc-2024/internal"
	"aoc-2024/internal/taoc"
	"cmp"
	"fmt"
	"github.com/alitto/pond/v2"
	"github.com/dominikbraun/graph"
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/sa-/slicefunk"
	"github.com/samber/lo"
	"github.com/samber/lo/parallel"
	"log/slog"
	"sort"
	"strings"
	"time"
)

const DemoInput = `
AAAA
BBCD
BBCC
EEEC
`
const DemoInputE = `
EEEEE
EXXXX
EEEEE
EXXXX
EEEEE
`

const DemoInput1 = `
OOOOO
OXOXO
OOOOO
OXOXO
OOOOO
`

// 00
// 0X

const DemoInput2 = `
RRRRIICCFF
RRRRIICCCF
VVRRRCCFFF
VVRCCCJFFF
VVVVCJJCFE
VVIVCCJJEE
VVIIICJJEE
MIIIIIJJEE
MIIISIJEEE
MMMISSJEEE
`

const DemoInput3 = `
AAAAAA
AAABBA
AAABBA
ABBAAA
ABBAAA
AAAAAA
`

var Input string

func init() {
	//Input = DemoInput2
	Input = internal.Download(2024, 12)
}

type Region struct {
	Main   taoc.Position
	Flavor rune
	Plants map[taoc.Position]struct{}
}

type State struct {
	Graph graph.Graph[taoc.Position, taoc.Position]
	End   taoc.Position

	PlantToRegion map[taoc.Position]*Region
	Garden        map[taoc.Position]rune
}

type Plant = lo.Entry[taoc.Position, rune]

//goland:noinspection Flavor,t
func main() {
	lines := slicefunk.Map(strings.Split(strings.TrimSpace(Input), "\n"), strings.TrimSpace)
	garden := lo.FromEntries(lo.FlatMap(lines, func(line string, y int) []Plant {
		return lo.Map([]rune(line), func(h rune, x int) Plant {
			return Plant{Key: taoc.Position{Coord: taoc.Coord{A: x, B: y}}, Value: h}
		})
	}))

	s := State{
		Graph: graph.New(func(t taoc.Position) taoc.Position {
			return t
		}),
		End:    taoc.Position{Coord: taoc.Coord{A: len(lines[0]) - 1, B: len(lines) - 1}},
		Garden: garden,
	}

	slog.Info("filling graph...")

	lo.ForEach(lo.Entries(garden), func(p Plant, _ int) {
		_ = s.Graph.AddVertex(p.Key)

		lo.ForEach(taoc.Directions4, func(d taoc.Vector, _ int) {
			adj := p.Key.Add(d)
			if !adj.In(s.End) {
				return
			}

			if p.Value != garden[adj] {
				return
			}

			_ = s.Graph.AddEdge(p.Key, adj)
		})
	})

	slog.Info("computing transitive closure...")

	TransitiveClosure(s)

	//enc := gob.NewEncoder(lo.Must(os.Create(path.Join(
	//	lo.Must(os.Getwd()),
	//	"graph.bin",
	//))))
	//
	//gob.Register(s.Graph)
	//lo.Must0(enc.Encode(s))

	slog.Info("computing regions...")

	adjMap := lo.Must(s.Graph.AdjacencyMap())
	posToRegion := make(map[taoc.Position]*Region)
	regions := make(map[taoc.Position]*Region)
	for v, adjs := range adjMap {
		r := posToRegion[v]
		if r == nil {
			r = &Region{Main: v, Flavor: s.Garden[v], Plants: map[taoc.Position]struct{}{v: {}}}
			regions[r.Main] = r
			posToRegion[v] = r
		}
		for adj := range adjs {
			r.Plants[adj] = struct{}{}
			posToRegion[adj] = r
		}
	}

	slog.Info("computing price...")

	price := lo.SumBy(lo.Values(regions), func(r *Region) int {
		volume := len(r.Plants)

		perimeter := lo.SumBy(lo.Keys(r.Plants), func(p taoc.Position) int {
			return lo.SumBy(taoc.Directions4, func(d taoc.Vector) int {
				adj := p.Add(d)
				if !adj.In(s.End) {
					return 1
				}

				if r.Main == posToRegion[adj].Main {
					return 0
				}

				return 1
			})
		})

		slog.Info(
			"region",
			"Flavor", string(s.Garden[r.Main]),
			"volume", volume,
			"perimeter", perimeter,
			"price", volume*perimeter,
		)

		return volume * perimeter
	})

	fmt.Println(price)
	slog.Info("now by slots")

	// thanks reddit

	// compute number of corners for each region
	price = lo.Sum(parallel.Map(lo.Values(regions), func(r *Region, _ int) int {
		corners := lo.SumBy(lo.Keys(r.Plants), func(p taoc.Position) int {
			my := s.Garden[p]
			return lo.SumBy(taoc.Corners, func(c [2]taoc.Vector) int {
				one := s.Garden[p.Add(c[0])]
				two := s.Garden[p.Add(c[1])]

				if my != one && my != two {
					return 1
				}
				if one == two && one != s.Garden[p.Add(c[0]).Add(c[1])] {
					// inner
					return 1
				}
				return 0
			})
		})

		slog.Info("region", "Flavor", string(s.Garden[r.Main]), "corners", corners, "Plants", len(r.Plants), "price", corners*len(r.Plants))

		return corners * len(r.Plants)
	}))

	fmt.Println(price)

	//file := lo.Must(os.Create("./simple.gv"))
	//lo.Must0(draw.DOT(s.Graph, file))
}

func TransitiveClosure(s State) {
	type k lo.Tuple2[taoc.Position, taoc.Position]

	tc := xsync.NewMapOf[k, struct{}]()

	p := pond.NewPool(0)

	g := p.NewGroup()

	t := lo.Keys(s.Garden)
	sort.Slice(t, func(i, j int) bool {
		return cmp.Or(
			cmp.Less(t[i].A, t[j].A),
			cmp.Less(t[i].B, t[j].B),
		)
	})

	c := xsync.NewCounter()
	lo.ForEach(t, func(v taoc.Position, _ int) {
		time.Sleep(10 * time.Millisecond)
		g.Submit(func() {
			lo.Must0(graph.DFS(s.Graph, v, func(n taoc.Position) bool {
				tc.Store(k{A: v, B: n}, struct{}{})
				return false
			}))

			c.Add(1)
			if c.Value()%500 == 0 {
				slog.Info("done computing TC for", "c", c.Value(), "of", len(t))
			}
		})
	})

	lo.Must0(g.Wait())

	tc.Range(func(k k, _ struct{}) bool {
		_ = s.Graph.AddEdge(k.A, k.B)
		return true
	})
}
