package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
)

type textdata struct {
	data []string
}

func (t *textdata) input(f io.Reader) {
	scanner := bufio.NewScanner(f)

	if t.data == nil {
		t.data = make([]string, 0)
	}

	// todo scan options
	for scanner.Scan() {
		someline := scanner.Text()
		t.data = append(t.data, someline)
	}
}

func (t *textdata) parse() {
	fmt.Println(len(t.data))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	t := new(textdata)

	t.input(os.Stdin)
	t.parse()
}
