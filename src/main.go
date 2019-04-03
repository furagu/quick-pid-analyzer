package main

import (
	"fmt"
	"log"

	"github.com/furagu/quick-pid-analyzer/src/bblog"
	"github.com/furagu/quick-pid-analyzer/src/trace"
)

func main() {
	fmt.Println("Loading a log")

	s, err := bblog.SampleSession()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parsing traces\n")

	r, p, y, err := trace.TracesFromLogSession(s)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range []*trace.Trace{r, p, y} {
		t.Print()
		fmt.Println("")
	}
}
