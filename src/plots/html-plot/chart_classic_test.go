package htmlPlot

import (
	"bytes"
	"fmt"
	"github.com/furagu/quick-pid-analyzer/src/plots"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestClassicChartViews(t *testing.T) {
	pointsFoo := []plots.Point{{0, 0}, {5, 10}, {10, 0}}
	pointsBar := []plots.Point{{0, 5}, {10, 5}}
	lineFoo := &plots.Line{
		Points: pointsFoo,
	}

	t.Run("should display axis names", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{
			XAxisName: "x-foo",
			YAxisName: "y-foo",
		})
		ch.AddLine(lineFoo)
		testSnapshotChart(t, "axis_names", ch)
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
		testSnapshotChart(t, "legend", ch)
	})

	t.Run("should display grid", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{
			ShowGrid: true,
		})
		ch.AddLine(lineFoo)
		testSnapshotChart(t, "grid", ch)
	})

	t.Run("should show custom line color", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(&plots.Line{
			Points:      pointsFoo,
			StrokeColor: &plots.Color{255, 0, 0, 255},
		})
		testSnapshotChart(t, "stroke_color", ch)
	})

	t.Run("should show custom line width", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(&plots.Line{
			Points:      pointsFoo,
			StrokeWidth: 5,
		})
		testSnapshotChart(t, "stroke_width", ch)
	})

	t.Run("should show custom filling color", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{})
		ch.AddLine(&plots.Line{
			Points:    pointsFoo,
			FillColor: &plots.Color{255, 0, 0, 50},
		})
		testSnapshotChart(t, "filling_color", ch)
	})

	t.Run("should apply custom viewBox", func(t *testing.T) {
		ch := NewClassicChart("", ChartConfig{
			ViewBox: &[4]int{50, 50, 500, 500},
		})
		ch.AddLine(lineFoo)
		testSnapshotChart(t, "view_box", ch)
	})
}

func testSnapshotChart(t *testing.T, testName string, chart plots.Chart) {
	buf := &bytes.Buffer{}
	err := chart.Draw(buf)
	assert.NoError(t, err)

	dir := "./_snapshot/chart_classic"
	fileName := fmt.Sprintf("%s/%s.svg", dir, testName)
	failedFileName := fmt.Sprintf("%s/%s.new.svg", dir, testName)

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
