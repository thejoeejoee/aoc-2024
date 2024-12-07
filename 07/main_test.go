package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestEquation(t *testing.T) {
	e := Equation{A: 33, B: []int64{3, 5, 6}}

	has, rest := e.resolve()

	assert.False(t, has)
	assert.Equal(t, []Equation{
		{30, []int64{5, 6}},
		{11, []int64{5, 6}},
	}, rest)

	eqs := []Equation{
		{A: 11, B: []int64{5, 6}},
		{A: 30, B: []int64{5, 6}},
	}

	for _, eq := range eqs {
		has, rest = eq.resolve()

		assert.True(t, has)
		assert.Empty(t, rest)
	}

	eqs = []Equation{
		{A: 12, B: []int64{5, 6}},
		{A: 30, B: []int64{5, 5}},
	}

	for _, eq := range eqs {
		has, rest = eq.resolve()

		assert.False(t, has)
		assert.Empty(t, rest)
	}

}

func TestConcatOp(t *testing.T) {
	cases := []struct {
		a, b int64
		want int64
		ok   bool
	}{
		{a: 123, b: 23, want: 1, ok: true},
		{a: 123, b: 3, want: 12, ok: true},
		{a: 123, b: 2, want: 0, ok: false},
		{a: 123, b: 1, want: 0, ok: false},
		{a: 123, b: 0, want: 0, ok: false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c), func(t *testing.T) {
			got, ok := concatOp(c.a, c.b)
			assert.Equal(t, c.ok, ok)
			assert.Equal(t, c.want, got)
		})
	}
}

var concatOp = func(result, suffix int64) (int64, bool) {
	// remove suffix from result by log 10
	rl := math.Ceil(math.Log10(float64(result)))
	sl := math.Ceil(math.Log10(float64(suffix)))

	if rl <= sl || result < 0 {
		return 0, false
	}

	base := int64(math.Pow(10, sl))

	if result%base != suffix {
		return 0, false
	}

	return result / base, true
}

func FuzzConcatMath(f *testing.F) {

	f.Add(123, 123)
	f.Add(123, 23)
	f.Add(123, 3)
	f.Add(123, 2)
	f.Add(123, 1)
	f.Add(123, 0)
	//f.Fuzz(func(t *testing.T, result int, suffix int) {
	//	concatOp(int64(result), int64(suffix))
	//})
	f.Fuzz(func(t *testing.T, result int, suffix int) {
		stringOp(int64(result), int64(suffix))
	})
}
