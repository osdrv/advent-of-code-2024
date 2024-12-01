package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const BENCH_RUNS = 10

func startsWith(s string, p string) bool {
	if len(s) < len(p) {
		return false
	}
	return s[:len(p)] == p
}

func getDayFrom(dir fs.FileInfo) int {
	day, err := strconv.Atoi(dir.Name()[4:])
	if err != nil {
		panic(err)
	}
	return day
}

func benchmark(dir fs.FileInfo, nRuns int, execCmd string, args ...string) []time.Duration {
	os.Chdir(dir.Name())
	defer os.Chdir("..")

	runs := make([]time.Duration, 0, nRuns)

	if err := exec.Command("make", "build").Run(); err != nil {
		panic(err)
	}

	for run := 0; run < nRuns; run++ {
		start := time.Now()
		cmd := exec.Command(execCmd, args...)
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		elapsed := time.Now().Sub(start)
		runs = append(runs, elapsed)
	}

	if err := exec.Command("make", "clean").Run(); err != nil {
		panic(err)
	}

	sort.Slice(runs, func(i, j int) bool {
		return runs[i].Nanoseconds() < runs[j].Nanoseconds()
	})

	return runs
}

func getLOC(dir fs.FileInfo, fNames ...string) int {
	files, err := ioutil.ReadDir(dir.Name())
	if err != nil {
		panic(err)
	}

	fMap := make(map[string]int, len(fNames))
	for _, fName := range fNames {
		fMap[fName] = 0
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if _, ok := fMap[file.Name()]; !ok {
			continue
		}
		data, err := os.ReadFile(dir.Name() + "/" + file.Name())
		if err != nil {
			panic(err)
		}
		loc := len(bytes.Split(data, []byte{'\n'}))
		fMap[file.Name()] = loc - 1
	}

	totalLoc := 0
	for _, loc := range fMap {
		totalLoc += loc
	}

	return totalLoc
}

func main() {

	includeDays := make(map[string]bool)
	idArg := flag.String("days", "", "Only include enlisted days")

	flag.Parse()
	for _, day := range strings.Split(*idArg, ",") {
		if len(day) == 0 {
			continue
		}
		includeDays[day] = true
	}

	dirs, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}
	dayDirs := make([]fs.FileInfo, 0, 1)
	for _, dir := range dirs {
		if !dir.IsDir() || !startsWith(dir.Name(), "day-") {
			continue
		}
		if len(includeDays) > 0 && !includeDays[dir.Name()] {
			continue
		}
		dayDirs = append(dayDirs, dir)
	}

	sort.Slice(dayDirs, func(i, j int) bool {
		return getDayFrom(dayDirs[i]) < getDayFrom(dayDirs[j])
	})

	for _, dir := range dayDirs {
		fmt.Printf("=== Benchmarking %s ===\n", dir.Name())
		elapsed := benchmark(dir, BENCH_RUNS, "make", "run-build")
		fmt.Printf("elapsed: %+v\n", elapsed)
		fmt.Printf(
			"p0: %d, p50: %d, p70: %d, p90: %d, p100: %d\n",
			elapsed[0].Milliseconds(),
			elapsed[int(float64(len(elapsed))*0.5)-1].Milliseconds(),
			elapsed[int(float64(len(elapsed))*0.7)-1].Milliseconds(),
			elapsed[int(float64(len(elapsed))*0.9)-1].Milliseconds(),
			elapsed[len(elapsed)-1].Milliseconds(),
		)
		fmt.Printf("Lines of Code: %d\n", getLOC(dir, "main.go"))
	}
}
