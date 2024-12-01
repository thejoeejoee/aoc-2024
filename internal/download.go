package internal

import (
	"fmt"
	"github.com/samber/lo"
	"io"
	"net/http"
	"os"
)

func Download(year, day int) string {
	cacheFile := fmt.Sprintf("cache/%d-%d.txt", year, day)

	if _, err := os.Stat(cacheFile); err == nil {
		bytes := lo.Must(os.ReadFile(cacheFile))
		return string(bytes)
	}

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	r := lo.Must(http.NewRequest("GET", url, nil))
	r.Header.Set("Cookie", fmt.Sprintf("session=%s", os.Getenv("AOC_SESSION")))
	resp := lo.Must(http.DefaultClient.Do(r))
	bytes := lo.Must(io.ReadAll(resp.Body))

	if resp.StatusCode != http.StatusOK {
		panic(string(bytes))
	}

	lo.Must0(os.WriteFile(cacheFile, bytes, 0644))

	return string(bytes)
}
