package internal

import (
	"aoc-2024/cache"
	"fmt"
	"github.com/samber/lo"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
)

func Download(year, day int) string {

	name := fmt.Sprintf("%d-%d.txt", year, day)
	dir := fmt.Sprintf("%s/cache", lo.Must(os.Getwd()))

	file, err := cache.FS.ReadFile(name)
	switch {
	case err == nil:
		slog.Info("using cached input", "file", name)
		return string(file)
	case !os.IsNotExist(err):
		panic(err)
	default:
	}

	full := path.Join(dir, name)

	slog.Info("downloading input", "year", year, "day", day, "file", full)

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	r := lo.Must(http.NewRequest("GET", url, nil))
	r.Header.Set("Cookie", fmt.Sprintf("session=%s", os.Getenv("AOC_SESSION")))
	resp := lo.Must(http.DefaultClient.Do(r))
	bytes := lo.Must(io.ReadAll(resp.Body))

	if resp.StatusCode != http.StatusOK {
		panic(string(bytes))
	}

	lo.Must0(os.WriteFile(full, bytes, 0644))

	return string(bytes)
}
