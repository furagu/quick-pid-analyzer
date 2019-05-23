package htmlPlot

import (
	"bytes"
	"github.com/furagu/quick-pid-analyzer/src/plots"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"io"
	"sync"
	"text/template"
)

var (
	// tplBox contains statically loaded templates from the "template" dir
	tplBox = packr.NewBox("./template")

	plotTplStr, _ = tplBox.FindString("plot.html")
	plotTpl       = template.Must(template.New("plot").Parse(plotTplStr))
)

type PlotConfig struct {
	// How many charts per row to display
	Cols int
}

// HtmlPlotter implements plots.HtmlPlotter interface, and plots multiple SVG plots to
// a single HTML file.
type HtmlPlotter struct {
	cfg    PlotConfig
	charts []plots.Chart
	wg     sync.WaitGroup
}

// NewHtmlPlotter crates a new plotter instance.
func NewHtmlPlotter(cfg PlotConfig) *HtmlPlotter {
	return &HtmlPlotter{
		cfg: cfg,
	}
}

// AddChart adds another SVG chart to HTML plot
func (p *HtmlPlotter) AddChart(c plots.Chart) {
	p.charts = append(p.charts, c)
}

// Plot generates HTML file with SVG charts and writes it to the target.
func (p *HtmlPlotter) Plot(target io.Writer) error {
	var charts []io.Reader
	for _, chart := range p.charts {
		buf := &bytes.Buffer{}
		if err := p.plotChart(chart, buf); err != nil {
			return err
		}
		charts = append(charts, buf)
	}

	cols := p.cfg.Cols
	if cols == 0 {
		cols = 2
	}

	return plotTpl.Execute(target, map[string]interface{}{
		"ChartWidthPercents": 100 / cols,
		"Charts":             charts,
	})
}

func (p *HtmlPlotter) plotChart(chart plots.Chart, target io.Writer) error {
	err := chart.Draw(target)
	return errors.Wrapf(err, "could not draw chart '%s' on the plot", chart.Name())
}
