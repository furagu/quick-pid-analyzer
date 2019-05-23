package htmlPlot

import (
	"bytes"
	"fmt"
	"github.com/furagu/quick-pid-analyzer/src/plots"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"testing"
)

func TestHtmlPlotter_Plot(t *testing.T) {
	pointsFoo := []plots.Point{{0, 0}, {5, 10}, {10, 0}}
	pointsBar := []plots.Point{{0, 5}, {10, 5}}
	lineFoo := &plots.Line{
		Points: pointsFoo,
	}

	t.Run("should show custom amount cols per row", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(lineFoo)
		testSnapshotPlotWithCols(t, "cols", 3, ch)
	})

	t.Run("should display axis names", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{
			XAxisName: "x-foo",
			YAxisName: "y-foo",
		})
		ch.AddLine(lineFoo)
		testSnapshotPlot(t, "axis_names", ch)
	})

	t.Run("should show legend for each line", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{
			ShowLegend: true,
		})
		ch.AddLine(&plots.Line{
			Name:   "foo",
			Points: pointsFoo,
		})
		ch.AddLine(&plots.Line{
			Name:   "bar",
			Points: pointsBar,
		})
		testSnapshotPlot(t, "legend", ch)
	})

	t.Run("should display grid", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{
			ShowGrid: true,
		})
		ch.AddLine(lineFoo)
		testSnapshotPlot(t, "grid", ch)
	})

	t.Run("should show custom line color", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(&plots.Line{
			Points:      pointsFoo,
			StrokeColor: &plots.Color{255, 0, 0, 255},
		})
		testSnapshotPlot(t, "stroke_color", ch)
	})

	t.Run("should show custom line width", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(&plots.Line{
			Points:      pointsFoo,
			StrokeWidth: 5,
		})
		testSnapshotPlot(t, "stroke_width", ch)
	})

	t.Run("should show custom filling color", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(&plots.Line{
			Points:    pointsFoo,
			FillColor: &plots.Color{255, 0, 0, 50},
		})
		testSnapshotPlot(t, "filling_color", ch)
	})

	t.Run("should apply custom viewBox", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{
			ViewBox: &[4]int{50, 50, 500, 500},
		})
		ch.AddLine(lineFoo)
		testSnapshotPlot(t, "view_box", ch)
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

func testSnapshotPlot(t *testing.T, testName string, charts ...plots.Chart) {
	testSnapshotPlotWithCols(t, testName, 0, charts...)
}

func testSnapshotPlotWithCols(t *testing.T, testName string, cols int, charts ...plots.Chart) {
	pl := NewHtmlPlotter(PlotConfig{Cols: cols})
	for _, chart := range charts {
		pl.AddChart(chart)
	}

	buf := &bytes.Buffer{}
	err := pl.Plot(buf)
	assert.NoError(t, err)

	fileName := fmt.Sprintf("./_snapshot/%s.html", testName)
	failedFileName := fmt.Sprintf("./_snapshot/%s.new.html", testName)

	if currentSnapshot, err := ioutil.ReadFile(fileName); err == nil {
		if !assert.Equal(t, currentSnapshot, buf.Bytes()) {
			err = ioutil.WriteFile(failedFileName, buf.Bytes(), 0644)
		}

	} else {
		// save snapshot if not saved yet
		err = ioutil.WriteFile(fileName, buf.Bytes(), 0644)
		assert.NoError(t, err)
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
