package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

type textdata struct {
	rawdata []string
	data    []string
}

func (t *textdata) input(f io.Reader) {
	scanner := bufio.NewScanner(f)

	if t.data == nil {
		t.rawdata = make([]string, 0)
	}

	for scanner.Scan() {
		someline := scanner.Text()
		t.rawdata = append(t.rawdata, someline)
	}
}

func (t *textdata) tf(term string, norm bool) float64 {
	return 0.0
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func parse(lines []string) <-chan string {
	out := make(chan string, 8*100)

	go func() {
		for _, line := range lines {
			if len(line) > 0 {
				out <- line
			}
		}

		close(out)
	}()

	return out
}

func parselower(in <-chan string) <-chan string {
	out := make(chan string, cap(in))

	go func() {
		for line := range in {
			out <- strings.ToLower(line)
		}
		close(out)
	}()

	return out
}

func parsesplit(in <-chan string) <-chan string {
	out := make(chan string, cap(in))

	go func() {
		for line := range in {
			grams := strings.Split(line, " ")
			var prev string
			for _, gram := range grams {
				if len(gram) > 0 {
					out <- (prev + gram)
				}
				prev = gram
			}
		}
		close(out)
	}()

	return out
}

func resulter(c map[string]int, in <-chan string) int {
	var total int

	for gram := range in {
		total += 1
		c[gram] += 1
	}

	return total
}

func lexd(total, unique int) int {
	return int(float64(unique) / float64(total) * 100)
}

func main() {
	t := new(textdata)

	t.input(os.Stdin)

	counts := make(map[string]int)

	parser := parse(t.rawdata)
	parselower := parselower(parser)
	parsespliter := parsesplit(parselower)
	total := resulter(counts, parsespliter)

	fmt.Println(lexd(total, len(counts)))
}
