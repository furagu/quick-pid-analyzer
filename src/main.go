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

		plots.StepResponseFromTrace(t)
		p := plots.StepResponseFromTrace(t)
		low, high := p.Calculate()

		fmt.Println("Low\n")
		for _, x := range low {
			fmt.Printf("%.12f\t%.12f\n", x[0], x[1])
		}

		fmt.Println("High\n")
		for _, x := range high {
			fmt.Printf("%.12f %.12f\n", x[0], x[1])
		}
	}
}
