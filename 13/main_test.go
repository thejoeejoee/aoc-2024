package main

import (
	"cmp"
	"fmt"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

var sortResults = func(s []lo.Tuple2[int, int]) func(i, j int) bool {
	return func(i, j int) bool {
		return cmp.Or(
			cmp.Less(s[i].A, s[j].A),
			cmp.Less(s[i].B, s[j].B),
		)

	}
}

func Test_allSums(t *testing.T) {
	tests := []struct {
		target int
		a      int
		b      int
		want   []lo.Tuple2[int, int]
	}{
		{3, 1, 2, []lo.Tuple2[int, int]{
			{1, 1}, {3, 0},
		}},
		{5, 1, 2, []lo.Tuple2[int, int]{
			{1, 2}, {3, 1}, {5, 0},
		}},
		{6, 1, 2, []lo.Tuple2[int, int]{
			{0, 3}, {2, 2}, {4, 1}, {6, 0},
		}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.target), func(t *testing.T) {
			got, err := allSums(tt.target, tt.a, tt.b)

			assert.NoError(t, err)

			sort.Slice(got, sortResults(got))
			sort.Slice(tt.want, sortResults(tt.want))

			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
