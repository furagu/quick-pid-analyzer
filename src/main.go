package main

import (
	"fmt"
	"log"

	"github.com/furagu/quick-pid-analyzer/src/bblog"
	"github.com/furagu/quick-pid-analyzer/src/plots"
	"github.com/furagu/quick-pid-analyzer/src/trace"
)

func main() {
	fmt.Println("Loading log data...")

	s, err := bblog.SampleSession()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d samples.\n", len(s.Values))

	fmt.Println("\nParsing traces...\n")

	r, p, y, err := trace.TracesFromLogSession(s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Analyzing...\n")

	for _, t := range []*trace.Trace{r, p, y} {
		t.Print()
		fmt.Println("")

		p := plots.StepResponseFromTrace(t)
		x, y := p.Calculate()

		fmt.Println(x)
		fmt.Println(y)

		fmt.Println("")
	}
}
