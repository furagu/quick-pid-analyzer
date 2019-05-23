package htmlPlot

import (
	"bytes"
	"fmt"
	"github.com/furagu/quick-pid-analyzer/src/plots"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestPlotter_Plot(t *testing.T) {
	pointsFoo := []plots.Point{{0, 0}, {5, 10}, {10, 0}}
	pointsBar := []plots.Point{{0, 5}, {10, 5}}
	lineFoo := plots.Line{
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
		ch.AddLine(plots.Line{
			Name:   "foo",
			Points: pointsFoo,
		})
		ch.AddLine(plots.Line{
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
		ch.AddLine(plots.Line{
			Points:      pointsFoo,
			StrokeColor: &plots.Color{255, 0, 0, 255},
		})
		testSnapshotPlot(t, "stroke_color", ch)
	})

	t.Run("should show custom line width", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(plots.Line{
			Points:      pointsFoo,
			StrokeWidth: 5,
		})
		testSnapshotPlot(t, "stroke_width", ch)
	})

	t.Run("should show custom filling color", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(plots.Line{
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
