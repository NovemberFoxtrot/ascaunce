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
	freq    map[string]int
	total   int
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

func (t *textdata) resulter(in <-chan string) {
	if t.freq == nil {
		t.freq = make(map[string]int)
	}

	for gram := range in {
		t.total += 1
		t.freq[gram] += 1
	}
}

func (t *textdata) lexd() float64 {
	return float64(len(t.freq)) / float64(t.total)
}

func (t *textdata) tf(term string) float64 {
	return float64(t.freq[term]) / float64(t.total)
}

func main() {
	for _, args := range os.Args[1:] {
		f, err := os.Open(args)
		t := new(textdata)

		if err != nil {
			fmt.Println(err)
		}

		defer f.Close()

		t.input(f)

		parser := parse(t.rawdata)
		parselower := parselower(parser)
		parsespliter := parsesplit(parselower)
		t.resulter(parsespliter)

		fmt.Println(t.lexd(), t.tf("the"))
	}
}
