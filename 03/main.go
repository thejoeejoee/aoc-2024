package main

import (
	"aoc-2024/internal"
	"fmt"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

const DemoInput = `xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`
const DemoInput2 = `xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`

var Input string = DemoInput2

func init() {
	Input = internal.Download(2024, 3)
	//Input = DemoInput2

}

var re = regexp.MustCompile(`mul\(\d+,\d+\)`)
var enablerRe = regexp.MustCompile(`don't\(\)|do\(\)|mul\(\d+,\d+\)`)

func main() {
	all := re.FindAllString(strings.TrimSpace(Input), -1)
	sum := lo.SumBy(all, func(m string) int {
		var a, b int
		_ = lo.Must(fmt.Sscanf(m, "mul(%d,%d)", &a, &b))
		return a * b
	})
	fmt.Println(sum)

	all = enablerRe.FindAllString(strings.TrimSpace(Input), -1)
	enabled := 1
	sum = lo.SumBy(all, func(m string) int {
		switch m[:3] {
		case "don":
			enabled = 0
			return 0
		case "do(":
			enabled = 1
			return 0
		default:
		}

		var a, b int
		_ = lo.Must(fmt.Sscanf(m, "mul(%d,%d)", &a, &b))
		return enabled * a * b
	})
	fmt.Println(sum)
}
