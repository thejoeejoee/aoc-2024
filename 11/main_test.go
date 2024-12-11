package main

import (
	"github.com/samber/lo"
	"strings"
	"testing"
)

func BenchmarkCompute(b *testing.B) {
	s := state{counts: lo.FromEntries(lo.Map(strings.Fields(strings.TrimSpace(Input)), func(s string, _ int) lo.Entry[string, uint64] {
		return lo.Entry[string, uint64]{
			Key:   s,
			Value: 1,
		}
	}))}

	for i := 0; i < b.N*1000; i++ {
		s = s.step()
	}
}
