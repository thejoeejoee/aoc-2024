package taoc

import "github.com/samber/lo"

var Top = Vector{Coord{A: 0, B: -1}}
var Right = Vector{Coord{A: 1, B: 0}}
var Bottom = Vector{Coord{A: 0, B: 1}}
var Left = Vector{Coord{A: -1, B: 0}}
var TopLeft = Vector{Coord{A: -1, B: -1}}
var TopRight = Vector{Coord{A: 1, B: -1}}
var BottomLeft = Vector{Coord{A: -1, B: 1}}
var BottomRight = Vector{Coord{A: 1, B: 1}}

var Directions4 = []Vector{Top, Right, Bottom, Left}

var Corners = [][2]Vector{
	{Top, Left},
	{Top, Right},
	{Bottom, Left},
	{Bottom, Right},
}

type Coord lo.Tuple2[int, int]

type Position struct {
	Coord
}

type Vector struct {
	Coord
}

func (c Position) Add(o Vector) Position {
	return Position{Coord: Coord{
		A: c.A + o.A,
		B: c.B + o.B,
	}}
}

func (c Position) Diff(o Position) Vector {
	return Vector{Coord: Coord{
		A: c.A - o.A,
		B: c.B - o.B,
	}}
}

func (c Position) In(end Position) bool {
	return c.A >= 0 && c.A <= end.A && c.B >= 0 && c.B <= end.B
}
