package htmlPlot

import (
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// NewClassicChart creates a preconfigured chart, which provides a classic view
// with enabled grid, axis names, legend, padding, etc.
func NewClassicChart(name string, cfg ChartConfig) *SvgChart {
	ch := &SvgChart{
		cfg:  cfg,
		name: name,
		baseLineStyle: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
		},
		graph: chart.Chart{
			XAxis: chart.XAxis{
				Name: cfg.XAxisName,
				Style: chart.Style{
					Show:     true,
					FontSize: 12,
				},
				NameStyle: chart.Style{
					Show:     true,
					FontSize: 14,
				},
			},
			YAxis: chart.YAxis{
				Name: cfg.YAxisName,
				NameStyle: chart.Style{
					Show:                true,
					TextRotationDegrees: 270,
					FontSize:            14,
				},
				Style:    chart.StyleShow(),
				AxisType: chart.YAxisSecondary,
			},
			Background: chart.Style{
				Padding: chart.Box{
					Top:  20,
					Left: 20,
				},
			},
			Series: []chart.Series{},
		},
	}

	if cfg.ShowLegend {
		ch.graph.Elements = []chart.Renderable{
			chart.Legend(&ch.graph, chart.Style{
				FontSize: 16,
			}),
		}
	}

	if cfg.ShowGrid {
		ch.graph.XAxis.GridMajorStyle = chart.Style{
			Show:        true,
			StrokeColor: drawing.Color{R: 0, G: 0, B: 0, A: 80},
			StrokeWidth: 1.0,
		}
		ch.graph.XAxis.GridMinorStyle = chart.Style{
			Show:        true,
			StrokeColor: drawing.Color{R: 0, G: 0, B: 0, A: 30},
			StrokeWidth: 1.0,
		}
		ch.graph.YAxis.GridMajorStyle = ch.graph.XAxis.GridMajorStyle
		ch.graph.YAxis.GridMinorStyle = ch.graph.XAxis.GridMinorStyle
	}

	return ch
}
