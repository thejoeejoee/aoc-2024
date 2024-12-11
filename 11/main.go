package main

import (
	"aoc-2024/internal"
	_ "aoc-2024/internal"
	"fmt"
	"github.com/samber/lo"

	"strconv"
	"strings"
)

const DemoInput = `125 17`

var Input string

func init() {
	//Input = DemoInput
	Input = internal.Download(2024, 11)
}

type state struct {
	counts map[string]uint64
}

func (old *state) String() string {
	return strings.Join(lo.FilterMap(lo.Entries(old.counts), func(item lo.Entry[string, uint64], _ int) (string, bool) {
		if item.Value == 0 {
			return "", false
		}
		return fmt.Sprintf("%dx'%s'", item.Value, item.Key), true
	}), " ")
}

func split(key string) (string, string) {
	l := key[:len(key)/2]
	r := lo.CoalesceOrEmpty(strings.TrimLeft(key[len(key)/2:], "0"), "0")
	//slog.Info("split", "key", key, "l", l, "r", r)
	return l, r
}

func (old *state) step() state {
	ns := state{counts: map[string]uint64{}}

	//If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
	//If the stone is engraved with a number that has an even number of digits, it is replaced by two counts.
	//The left half of the digits are engraved on the new left stone, and the right half of the digits are engraved on the new right stone.
	// (The new numbers don't keep extra leading zeroes: 1000 would become counts 10 and 0.)

	//If none of the other rules apply, the stone is replaced by a new stone; the old stone's number multiplied by 2024 is engraved on the new stone.

	ns.counts["1"] = old.counts["0"]

	evenKeys := lo.Filter(lo.Keys(old.counts), func(key string, _ int) bool {
		return 0 == (len(key) % 2)
	})
	oddKeys := lo.Filter(lo.Keys(old.counts), func(key string, _ int) bool {
		return 1 == (len(key)%2) && key != "0"
	})

	lo.ForEach(evenKeys, func(key string, _ int) {
		l, r := split(key)
		ns.counts[l] += old.counts[key]
		ns.counts[r] += old.counts[key]
	})
	lo.ForEach(oddKeys, func(key string, _ int) {
		k := strconv.Itoa(2024 * lo.Must(strconv.Atoi(key)))

		ns.counts[k] += old.counts[key]
	})
	return ns
}

func main() {
	s := state{counts: lo.FromEntries(lo.Map(strings.Fields(strings.TrimSpace(Input)), func(s string, _ int) lo.Entry[string, uint64] {
		return lo.Entry[string, uint64]{
			Key:   s,
			Value: 1,
		}
	}))}

	fmt.Println(s.String())
	for range 75 {
		s = s.step()
		fmt.Println(s.String())
	}

	o := lo.Sum(lo.Values(s.counts))

	fmt.Println(o)
}
