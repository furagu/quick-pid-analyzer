package htmlPlot

import (
	"github.com/furagu/quick-pid-analyzer/src/plots"
	"github.com/wcharczuk/go-chart/drawing"
	"io"

	"github.com/wcharczuk/go-chart"
)

// This viewBox argument provides best fit to make SVG tags responsive
var defaultViewBox = [4]int{0, 0, 1000, 400}

type ChartConfig struct {
	XAxisName  string
	YAxisName  string
	ShowLegend bool
	ShowGrid   bool
	ViewBox    *[4]int
}

// SvgChart is a plots.LineChart implementation, which exclusively prints SVG charts.
type SvgChart struct {
	cfg           ChartConfig
	name          string
	graph         chart.Chart
	baseLineStyle chart.Style
}

// Name returns the chart name (title)
func (c *SvgChart) Name() string {
	return c.name
}

// AddLine add a new line to the chart
func (c *SvgChart) AddLine(l plots.Line) {
	line := chart.ContinuousSeries{
		Name:            l.Name,
		Style:           c.getLineStyle(len(c.graph.Series), l),
		XValueFormatter: chart.ValueFormatter(l.XFormatter),
		YValueFormatter: chart.ValueFormatter(l.YFormatter),
	}
	for _, point := range l.Points {
		line.XValues = append(line.XValues, point[0])
		line.YValues = append(line.YValues, point[1])
	}

	c.graph.Series = append(c.graph.Series, line)
}

// Draw prints the chart to the target.
func (c *SvgChart) Draw(target io.Writer) error {
	// we need to wrap the original writer with our proxy,
	// which appends viewBox attribute to the <svg> tag,
	// otherwise SVG DOM element is not responsive
	viewBox := defaultViewBox
	if c.cfg.ViewBox != nil {
		viewBox = *c.cfg.ViewBox
	}
	target = &viewBoxWriter{
		Target:  target,
		ViewBox: viewBox,
	}

	return c.graph.Render(chart.SVG, target)
}

// resolveLineStyle creates a style configs for the line based on the baseStyle
// and the line's preferences
func (c SvgChart) getLineStyle(lineIdx int, line plots.Line) chart.Style {
	style := c.baseLineStyle

	color := plots.Color{0, 0, 0, 255}
	if line.StrokeColor != nil {
		color = *line.StrokeColor
	}
	style.StrokeColor = toDrawingColor(color)

	if line.FillColor != nil {
		style.FillColor = toDrawingColor(*line.FillColor)
	}

	width := line.StrokeWidth
	if width == 0 {
		width = 1
	}
	style.StrokeWidth = width

	return style
}

// toDrawingColor converts internal Color type to Color from the "chart" external library.
func toDrawingColor(c plots.Color) drawing.Color {
	return drawing.Color{R: c[0], G: c[1], B: c[2], A: c[3]}
}
