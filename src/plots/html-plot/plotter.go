package htmlPlot

import (
	"bytes"
	"github.com/furagu/quick-pid-analyzer/src/plots"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"io"
	"strings"
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
	cfg       PlotConfig
	charts    []plots.Chart
	wg        sync.WaitGroup
	latestErr error
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

	wg := &sync.WaitGroup{}
	errs := &errList{}
	for _, chart := range p.charts {
		buf := &bytes.Buffer{}
		charts = append(charts, buf)
		p.plotChartAsync(chart, buf, errs, wg)
	}
	wg.Wait()

	if !errs.empty() {
		return errs
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

func (p *HtmlPlotter) plotChartAsync(chart plots.Chart, target io.Writer, errs *errList, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := p.plotChart(chart, target); err != nil {
			errs.add(err)
		}
	}()
}

type errList struct {
	list []error
	mu   sync.Mutex
}

func (l *errList) add(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.list = append(l.list, err)
}

func (l *errList) empty() bool {
	return len(l.list) == 0
}

func (l errList) Error() string {
	var messages []string
	for _, err := range l.list {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "\n")
}
