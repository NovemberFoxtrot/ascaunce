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

	// todo scan options
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
			for _, gram := range grams {
				if len(gram) > 0 {
				  out <- gram
				}
			}
		}
		close(out)
	}()

	return out
}

func resulter(in <-chan string) {
	for gram := range in {
		fmt.Println(gram, len(gram))
	}
}

func main() {
	t := new(textdata)

	t.input(os.Stdin)

	parser := parse(t.rawdata)
	parselower := parselower(parser)
	parsespliter := parsesplit(parselower)
	resulter(parsespliter)
}
