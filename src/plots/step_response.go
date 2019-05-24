package plots

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"

	"github.com/furagu/quick-pid-analyzer/src/trace"
	"github.com/furagu/quick-pid-analyzer/src/utils"
	// "github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
	"github.com/mjibson/go-dsp/window"
	"gonum.org/v1/gonum/floats"
)

type Trace = trace.Trace

const (
	framelen  = 1.   // length of each single frame over which to compute response
	resplen   = 0.5  // length of respose window
	cutfreq   = 25.  // cut freqency of what is considered as input
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

func (s *StepResponse) Calculate() ([][2]float64, [][2]float64) {
	time := uniform(s.t.Time)

	input := interpolate(s.t.Time, calcInput(s.t), time)
	gyro := interpolate(s.t.Time, s.t.Gyro, time)

	flen := samplesInDuration(time, framelen)
	rlen := samplesInDuration(time, resplen)

	inputFrames := frames(input, flen, superpos)
	gyroFrames := frames(gyro, flen, superpos)

	win := window.Hann(flen)
	fmt.Printf("W: %s\n", utils.FloatsToString(win))

	for _, v := range inputFrames {
		floats.Mul(v, win)
	}
	for _, v := range gyroFrames {
		floats.Mul(v, win)
	}

	input_abs := absEach(inputFrames)
	input_avg := meanEach(input_abs)
	input_max := maxEach(input_abs)

	fmt.Printf("Input avg: %s\n", utils.FloatsToString(input_avg))
	fmt.Printf("Input max: %s\n", utils.FloatsToString(input_max))

	dt := time[0] - time[1]
	response := stack_response(inputFrames, gyroFrames, dt, rlen)

	low_mask, high_mask := lowHighMask(input_max, threshold)
	_, toolow_mask := lowHighMask(input_max, 20)

	da(high_mask)

	resp_time := time[0:rlen]
	floats.AddConst(-time[0], resp_time)

	resp_low := weightedModeAvg(response, resp_time, applyMask(low_mask, toolow_mask), -1.5, 3.5, 1000)
	resp_high := weightedModeAvg(response, resp_time, applyMask(high_mask, toolow_mask), -1.5, 3.5, 1000)

	// daa(resp_time, resp_low)
	// daa(resp_time, resp_low)

	// for _, v := range resp_low {
	// 	fmt.Printf("%+.8f\n", v)
	// }

	// os.Exit(0)

	// for _, v := range low_mask {
	// 	fmt.Printf("%+.8f\n", v)
	// }

	// for _, v := range high_mask {
	// 	fmt.Printf("%+.8f\n", v)
	// }

	// fmt.Println(utils.Ppaa(response))

	// plotTime := make([]float64, rlen)
	// copy(plotTime, time[0:rlen])
	// floats.AddConst(-time[0], plotTime)

	return toPlotData(resp_time, resp_low), toPlotData(resp_time, resp_high)
}

func toPlotData(t, r []float64) [][2]float64 {
	out := make([][2]float64, len(r))
	for i := range t {
		out[i] = [2]float64{t[i], r[i]}
	}
	return out
}

func weightedModeAvg(values [][]float64, time, weights []float64, min, max float64, bins int) []float64 {
	times := repeatArray(time, len(values))

	hist2d := histogram2d(flatten(times), time[0], time[len(time)-1], len(time),
		flatten(values), min, max, bins,
		flatten(repeatArray(weights, len(values[0]))))
	hist2d = transpose(hist2d)

	if emptyAll(hist2d) {
		return make([]float64, len(time), len(time))
	}

	maxes := maxEach0(hist2d)
	for _, v := range hist2d {
		floats.Div(v, maxes)
	}
	for _, v := range hist2d {
		floats.Mul(v, v)
	}

	resp_y := make([]float64, bins)
	floats.Span(resp_y, min, max)
	pixelpos := spreadInto2D(resp_y, len(time))

	avr := weightedAvr2d(pixelpos, hist2d)
	//     weights=hist2d_sm * hist2d_sm)

	// fmt.Printf("%s\n", utils.Ppaa(pixelpos))
	// fmt.Printf("%s\n", len(pixelpos[0]))
	// os.Exit(0)

	// da(pixelpos[0])

	// if hist2d.sum():
	//     hist2d_sm = hist2d
	//     hist2d_sm /= np.max(hist2d_sm, 0)

	//     pixelpos = np.repeat(resp_y.reshape(len(resp_y), 1), len(times[0]), axis=1)
	//     avr = np.average(pixelpos, 0, weights=hist2d_sm * hist2d_sm)

	return avr
}

func weightedAvr2d(a [][]float64, weights [][]float64) []float64 {
	out := make([]float64, len(a[0]))

	for j := 0; j < len(a[0]); j++ {
		weight := 0.
		for i := 0; i < len(a); i++ {
			out[j] += a[i][j] * weights[i][j]
			weight += weights[i][j]
		}
		out[j] /= weight
	}

	return out
}

func emptyAll(a [][]float64) bool {
	for _, v := range a {
		for _, w := range v {
			if w > 0 {
				return false
			}
		}
	}
	return true
}

func spreadInto2D(a []float64, l int) [][]float64 {
	out := make([][]float64, len(a))
	for i, x := range a {
		out[i] = make([]float64, l)
		floats.AddConst(x, out[i])
	}
	return out
}

func daa(a, b []float64) {
	fmt.Printf("Len: %d %d\n", len(a), len(b))
	for i := range a {
		fmt.Printf("%.12f\t%.12f\n", a[i], b[i])
	}
	os.Exit(0)
}

func da(a []float64) {
	fmt.Printf("Len: %d\n", len(a))
	for _, v := range a {
		fmt.Printf("%+.12f\n", v)
	}
	os.Exit(0)
}

func d(v float64) {
	fmt.Printf("%+.12f\n", v)
	os.Exit(0)
}

func histogram2d(x []float64, xmin, xmax float64, xbins int, y []float64, ymin, ymax float64, ybins int, weights []float64) [][]float64 {
	bins := make([][]float64, xbins)
	for i := range bins {
		bins[i] = make([]float64, ybins)
	}

	xedges := make([]float64, xbins+1)
	floats.Span(xedges, xmin, xmax)

	yedges := make([]float64, ybins+1)
	floats.Span(yedges, ymin, ymax)

	for i := 0; i < len(x); i++ {
		xi := x[i]
		yi := y[i]

		if xi < xmin || xi > xmax || yi < ymin || yi > ymax {
			continue
		}

		dx := sort.SearchFloat64s(xedges[1:], xi)
		dy := sort.SearchFloat64s(yedges[1:], yi)

		bins[dx][dy] += weights[i]
	}

	return bins
}

func transpose(a [][]float64) [][]float64 {
	out := make([][]float64, len(a[0]), len(a[0]))
	for i, v := range a {
		for j, w := range v {
			if out[j] == nil {
				out[j] = make([]float64, len(a), len(a))
			}
			out[j][i] = w
		}
	}
	return out
}

func repeatArray(a []float64, n int) [][]float64 {
	out := make([][]float64, n, n)
	for i := 0; i < n; i++ {
		out[i] = make([]float64, len(a), len(a))
		copy(out[i], a)
	}
	return out
}

func flatten(a [][]float64) []float64 {
	var out []float64
	for _, v := range a {
		out = append(out, v...)
	}
	return out
}

// def weighted_mode_avr(self, values, weights, vertrange, vertbins):

func applyMask(a, b []float64) []float64 {
	out := make([]float64, len(a), len(a))
	floats.MulTo(out, a, b)
	return out
}

func lowHighMask(s []float64, threshold float64) ([]float64, []float64) {
	low := make([]float64, len(s), len(s))
	high := make([]float64, len(s), len(s))

	for i, v := range s {
		if v > threshold {
			low[i] = 0
			high[i] = 1
		} else {
			low[i] = 1
			high[i] = 0
		}
	}
	return low, high
}

func stack_response(input, output [][]float64, dt float64, rlen int) [][]float64 {
	deconvolved := wienerDeconvolution(input, output, cutfreq, dt) // HERE
	deconvolved = slice2D(deconvolved, rlen)
	return cumSumEach(deconvolved)
}

func cumSumEach(in [][]float64) [][]float64 {
	out := make([][]float64, len(in))
	for i, v := range in {
		w := make([]float64, len(v))
		floats.CumSum(w, v)
		out[i] = w
	}
	return out
}

func slice2D(in [][]float64, l int) [][]float64 {
	out := make([][]float64, len(in), len(in))
	for i, v := range in {
		out[i] = v[:l]
	}
	return out
}

func wienerDeconvolution(input, output [][]float64, cutfreq float64, dt float64) [][]float64 {
	pad := 1024 - (len(input[0]) % 1024) // padding to power of 2, increases transform speed
	padEach(input, pad)
	padEach(output, pad)

	H := FFTEach(input)
	G := FFTEach(output)

	freq := abs(FFTFreq(len(input[0]), dt))
	sn := to_mask(clip(freq, cutfreq-1e-9, cutfreq))

	// sn_inverted := make([]float64, len(sn), len(sn))
	// floats.ScaleTo(sn_inverted, -1, sn)
	// floats.AddConst(1, sn_inverted)
	// len_lpf := floats.Sum(sn_inverted)
	// // fil := gaussian_filter1d(sn, len_lpf/6., "reflect") // TODO: implement later

	floats.Scale(-10., sn)
	floats.AddConst(10.*(1.+1e-9), sn) // +1e-9 to prohibit 0/0 situations

	Hcon := conjEach(H)

	sn_flipped := make([]float64, len(sn), len(sn))
	floats.AddConst(1., sn_flipped)
	floats.Div(sn_flipped, sn)

	return eachToReal(iFFTEach(divEach(mulEach(G, Hcon), addToCols(mulEach(H, Hcon), sn_flipped))))
}

func eachToReal(in [][]complex128) [][]float64 {
	out := make([][]float64, len(in), len(in))
	for i, v := range in {
		out[i] = make([]float64, len(v), len(v))
		for j, w := range v {
			out[i][j] = real(w)
		}
	}
	return out
}

func addToCols(a [][]complex128, b []float64) [][]complex128 {
	out := make([][]complex128, len(a), len(a))
	for i, v := range a {
		out[i] = make([]complex128, len(v), len(v))
		for j, w := range v {
			out[i][j] = w + complex(b[j], 0.)
		}
	}
	return out
}

func mulEach(a, b [][]complex128) [][]complex128 {
	out := make([][]complex128, len(a), len(a))
	for i, v := range a {
		out[i] = make([]complex128, len(v), len(v))
		for j, w := range v {
			out[i][j] = w * b[i][j]
		}
	}
	return out
}

func divEach(a, b [][]complex128) [][]complex128 {
	out := make([][]complex128, len(a), len(a))
	for i, v := range a {
		out[i] = make([]complex128, len(v), len(v))
		for j, w := range v {
			out[i][j] = w / b[i][j]
		}
	}
	return out
}

func conjEach(a [][]complex128) [][]complex128 {
	out := make([][]complex128, len(a), len(a))
	for i, v := range a {
		out[i] = make([]complex128, len(v), len(v))
		for j, w := range v {
			out[i][j] = cmplx.Conj(w)
		}
	}
	return out
}

func gaussian_filter1d(a []float64, sigma float64, mode string) []float64 {
	truncate := 4.

	lw := int(truncate*sigma + 0.5)
	weights := gaussian_kernel1d(sigma, lw)
	floats.Reverse(weights)

	return correlate1d(a, weights, mode)
}

func correlate1d(input, weights []float64, mode string) []float64 {
	filter_size := len(weights)
	size1 := filter_size / 2
	size2 := filter_size - size1 - 1
	EPSILON := math.Nextafter(1.0, 2.0) - 1.0

	fmt.Printf("filter_size: %s\n", filter_size)
	fmt.Printf("size1: %s\n", size1)
	fmt.Printf("size2: %s\n", size2)

	symmetric := 0
	if filter_size%2 > 0 {
		symmetric = 1
		for i := 1; i <= filter_size/2; i++ {
			if math.Abs(weights[i+size1]-weights[size1-i]) > EPSILON {
				symmetric = 0
				break
			}
		}

		if symmetric == 0 {
			symmetric = -1
			for i := 1; i <= filter_size/2; i++ {
				if math.Abs(weights[size1+i]+weights[size1-i]) > EPSILON {
					symmetric = 0
					break
				}
			}
		}
	}

	fmt.Printf("symmetric: %s\n", symmetric)

	return input

}

func gaussian_kernel1d(sigma float64, radius int) []float64 {
	p := []float64{0, 0, -0.5 / (sigma * sigma)}
	x := floatRange(-radius, radius)

	phi_exp := expEach(polynomialVal(p, x))
	floats.Scale(1/floats.Sum(phi_exp), phi_exp)

	return phi_exp
}

func floatRange(min, max int) []float64 {
	a := make([]float64, max-min+1)
	for i := range a {
		a[i] = float64(min + i)
	}
	return a
}

// Calculate the value of a polynomial with coefficients c at points x
func polynomialVal(c, x []float64) []float64 {
	out := make([]float64, len(x), len(x))

	floats.AddConst(c[len(c)-1], out)
	for i := len(c) - 2; i >= 0; i-- {
		floats.Mul(out, x)
		floats.AddConst(c[i], out)
	}
	return out
}

func expEach(a []float64) []float64 {
	out := make([]float64, len(a))
	for i := range a {
		out[i] = math.Exp(a[i])
	}
	return out
}

func clip(a []float64, min, max float64) []float64 {
	out := make([]float64, len(a), len(a))
	for i, v := range a {
		if v < min {
			v = min
		} else if v > max {
			v = max
		}
		out[i] = v
	}
	return out
}

func to_mask(a []float64) []float64 {
	out := make([]float64, len(a), len(a))
	copy(out, a)
	floats.AddConst(-floats.Min(out), out)
	floats.Scale(1/floats.Max(out), out)
	return out
}

func FFTFreq(n int, d float64) []float64 {
	out := make([]float64, n, n)
	N := int((float64(n)-1)/2 + 1)
	for i := 0; i < N; i++ {
		out[i] = float64(i)
	}
	for i := 0; i < n-N; i++ {
		out[N+i] = -float64(n)/2. + float64(i)
	}
	val := 1.0 / (float64(n) * d)
	floats.Scale(val, out)
	return out
}

func FFTEach(in [][]float64) [][]complex128 {
	out := make([][]complex128, len(in), len(in))
	for i, v := range in {
		out[i] = fft.FFT(toComplex(v))
	}
	return out
}

func iFFTEach(in [][]complex128) [][]complex128 {
	out := make([][]complex128, len(in), len(in))
	for i, v := range in {
		out[i] = fft.IFFT(v)
	}
	return out
}

func toComplex(in []float64) []complex128 {
	out := make([]complex128, len(in), len(in))
	for i, v := range in {
		out[i] = complex(v, 0.)
	}
	return out
}

func padEach(dst [][]float64, n int) {
	pad := make([]float64, n, n)
	for i, v := range dst {
		dst[i] = append(v, pad...)
	}
}

func maxEach0(a [][]float64) []float64 {
	out := make([]float64, len(a[0]), len(a[0]))
	for j := 0; j < len(a[0]); j++ {
		max := 0.
		for i := 0; i < len(a); i++ {
			if a[i][j] > max {
				max = a[i][j]
			}
		}
		out[j] = max
	}
	return out
}

func maxEach(a [][]float64) []float64 {
	out := make([]float64, len(a), len(a))
	for i, v := range a {
		out[i] = floats.Max(v)
	}
	return out
}

func meanEach(a [][]float64) []float64 {
	out := make([]float64, len(a), len(a))
	for i, v := range a {
		out[i] = floats.Sum(v) / float64(len(v))
	}
	return out
}

func absEach(a [][]float64) [][]float64 {
	out := make([][]float64, len(a), len(a))
	for i, v := range a {
		out[i] = abs(v)
	}
	return out
}

func abs(a []float64) []float64 {
	out := make([]float64, len(a), len(a))
	for i, v := range a {
		out[i] = math.Abs(v)
	}
	return out
}

func frames(samples []float64, length, subframes int) [][]float64 {
	step := int(length / subframes)
	n := int(len(samples)/step) - subframes
	frames := make([][]float64, n, n)
	for i := 0; i < n; i++ {
		frames[i] = make([]float64, length)
		copy(frames[i], samples[i*step:i*step+length])
	}
	return frames
}

func calcInput(t *Trace) []float64 {
	input := make([]float64, t.Len, t.Len)
	p_descale_factor := 1. / (0.032029 * t.PTerm) // 0.032029 is P scaling factor from Betaflight
	return floats.AddScaledTo(input, t.Gyro, p_descale_factor, t.P)
}

func samplesInDuration(time []float64, duration float64) int {
	step := time[1] - time[0]
	return int(duration / step)
}

func uniform(samples []float64) []float64 {
	resampled := make([]float64, len(samples), len(samples))
	return floats.Span(resampled, samples[0], samples[len(samples)-1])
}

func interpolate(x, y, newX []float64) []float64 {
	size := len(y)
	newY := make([]float64, size, size)

	// all inputs are assumed to be sorted
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
