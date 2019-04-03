package plots

import (
	"fmt"

	"github.com/furagu/quick-pid-analyzer/src/trace"
	"github.com/furagu/quick-pid-analyzer/src/utils"
	"gonum.org/v1/gonum/floats"
)

type Trace = trace.Trace

const (
	framelen  = 1.   // length of each single frame over which to compute response
	resplen   = 0.5  // length of respose window
	cutfreq   = 25.  // cutfreqency of what is considered as input
	superpos  = 16   // sub windowing (superpos windows in framelen)
	threshold = 500. // threshold for 'high input rate'
)

type StepResponse struct {
	t *Trace
}

func StepResponseFromTrace(trace *Trace) *StepResponse {
	return &StepResponse{
		t: trace,
	}
}

func (s *StepResponse) Calculate() (x, y []float64) {
	time := calcUniformTime(s.t)

	input := interpolate(s.t.Time, calcInput(s.t), time)
	throttle := interpolate(s.t.Time, s.t.Throttle, time)
	gyro := interpolate(s.t.Time, s.t.Gyro, time)

	fmt.Printf("Interpolated time: %s\n", utils.FloatsToString(time))
	fmt.Printf("Interpolated input: %s\n", utils.FloatsToString(input))
	fmt.Printf("Interpolated throttle: %s\n", utils.FloatsToString(throttle))
	fmt.Printf("Interpolated gyro: %s\n", utils.FloatsToString(gyro))

	return
}

func calcInput(t *Trace) []float64 {
	input := make([]float64, t.Len, t.Len)
	p_descale_factor := 1. / (0.032029 * t.PTerm) // 0.032029 is P scaling factor from Betaflight
	return floats.AddScaledTo(input, t.Gyro, p_descale_factor, t.P)
}

func calcUniformTime(t *Trace) []float64 {
	time := make([]float64, t.Len, t.Len)
	return floats.Span(time, t.Time[0], t.Time[t.Len-1])
}

func interpolate(x, y, newX []float64) []float64 { // Assumes all inputs are sorted
	size := len(y)
	newY := make([]float64, size, size)

	i := 0
	for j, nx := range newX {
		for {
			if x[i] > nx || i == size-1 {
				break
			}
			i++
		}

		var x1, y1, x2, y2 float64
		if i == 0 {
			x1 = x[i]
			y1 = y[i]
			x2 = x[i+1]
			y2 = y[i+1]
		} else {
			x1 = x[i-1]
			y1 = y[i-1]
			x2 = x[i]
			y2 = y[i]
		}

		slope := (y2 - y1) / (x2 - x1)
		ny := slope*(nx-x1) + y1
		newY[j] = ny
	}

	return newY
}
