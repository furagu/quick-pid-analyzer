package htmlPlot

import (
	"bytes"
	"fmt"
	"github.com/furagu/quick-pid-analyzer/src/plots"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestHtmlPlotter_Plot(t *testing.T) {
	lineFoo := &plots.Line{
		Points: []plots.Point{{0, 0}, {1, 1}},
	}

	chartFoo := NewClassicChart("foo", ChartConfig{})
	chartFoo.AddLine(lineFoo)

	chartBar := NewClassicChart("bar", ChartConfig{})
	chartBar.AddLine(lineFoo)

	t.Run("should render HTML file with multiple charts", func(t *testing.T) {
		pl := NewHtmlPlotter(PlotConfig{})
		pl.AddChart(chartFoo)
		pl.AddChart(chartBar)

		buf := &bytes.Buffer{}
		err := pl.Plot(buf)
		assert.NoError(t, err)
		assert.Regexp(t, `(?s)(class="chart".+?<svg\s+.+){2}`, buf.String(), "html should contain 2 chart items with svg tags")
	})

	t.Run("should show charts in 2 columns by default", func(t *testing.T) {
		pl := NewHtmlPlotter(PlotConfig{})
		pl.AddChart(chartFoo)

		buf := &bytes.Buffer{}
		err := pl.Plot(buf)
		assert.NoError(t, err)
		assert.Regexp(t, `(?s).chart\s*\{.+width:\s*50%;`, buf.String(), "every chart should be 50% wide to be displayed in 2 columns")
	})

	t.Run("should set custom amount of columns", func(t *testing.T) {
		pl := NewHtmlPlotter(PlotConfig{
			Cols: 3,
		})
		pl.AddChart(chartFoo)

		buf := &bytes.Buffer{}
		err := pl.Plot(buf)
		assert.NoError(t, err)
		assert.Regexp(t, `(?s).chart\s*\{.+width:\s*33%;`, buf.String(), "every chart should be 33% wide to be displayed in 3 columns")
	})
}

// every benchmark should consume the created result (store to a global var)
var benchResult []byte

// Previous benchmark results for testScale=10:
// With synchronous rendering: BenchmarkHtmlPlotter_Plot-8 20 76889518 ns/op
// With concurrent rendering:  BenchmarkHtmlPlotter_Plot-8 30 42574735 ns/op
func BenchmarkHtmlPlotter_Plot(b *testing.B) {
	testScale := 10 // how many plots with how many lines to generate

	// generate dataset
	var data [][][][2]float64
	for i := 0; i < testScale; i++ {
		var lines [][][2]float64
		for j := 0; j <= i; j++ {
			lines = append(lines, generateLine(lineOpts{
				xSize: 0.5,
				ySize: 2,
				xStep: 0.001,
				yStep: 0.008,
			}))
		}
		data = append(data, lines)
	}

	for i := 0; i < b.N; i++ {
		pl := NewHtmlPlotter(PlotConfig{})
		for j, lines := range data {
			chart := NewClassicChart(fmt.Sprintf("chart_%d", j), ChartConfig{
				ShowGrid:   true,
				ShowLegend: true,
				XAxisName:  "X values",
				YAxisName:  "Y values",
			})
			for k, points := range lines {
				chart.AddLine(&plots.Line{
					Name:        fmt.Sprintf("line_%d", k),
					StrokeWidth: 2,
					StrokeColor: &plots.Color{255, 0, 0, 255},
					FillColor:   &plots.Color{255, 0, 0, 50},
					Points:      points,
				})
			}
			pl.AddChart(chart)
		}

		target := &bytes.Buffer{}
		err := pl.Plot(target)
		assert.NoError(b, err)

		benchResult = target.Bytes()
	}
}

type lineOpts struct {
	xSize float64
	ySize float64
	xStep float64
	yStep float64
}

func generateLine(o lineOpts) [][2]float64 {
	prevLine := [2]float64{0, 0}
	lines := [][2]float64{prevLine}

	for x := o.xStep; x <= o.xSize; x += o.xStep {

		var direction float64 = 1
		if rand.Intn(2) == 0 {
			direction = -1
		}
		scale := float64(rand.Intn(3))

		y := prevLine[1] + (o.yStep * direction * scale)

		if y > o.ySize {
			y = o.ySize
		} else if y < 0 {
			y = 0
		}

		prevLine = [2]float64{x, y}
		lines = append(lines, prevLine)
	}

	return lines
}
