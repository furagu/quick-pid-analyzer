package plots

import "io"

// ValueFormatter is a function that takes a point value and produces a string.
type ValueFormatter func(v interface{}) string

// Point is a tuple of X,Y values
type Point = [2]float64

// Color represented by R,G,B bytes and last byte for Alpha-channel
type Color = [4]byte

// Plotter is responsible for plotting multiple charts to a single target
type Plotter interface {
	// Adds a new chart to the plot
	AddChart(Chart)
	// Plot plots all added charts and combines them in the single target,
	// (the strategy may differ, depends on the outcome format, e.g. png, svg, html)
	Plot(target io.Writer) error
}

// Chart represents a single chart with multiple items on it (lines, bars, circles, etc)
type Chart interface {
	// Name returns the chart name (title)
	Name() string
	// Draw generates a chart with provided lines and writes it to the target.
	Draw(target io.Writer) error
}

// LineChart is a chart which displays only lines
type LineChart interface {
	Chart
	// AddLine adds a new line to the chart
	AddLine(Line)
}

type Line struct {
	// Name is line name (used for legends, empty by default)
	Name string

	// Points is list of [x,y] points which represent the line
	Points []Point

	// StrokeWidth is the line thickness (1 by default)
	StrokeWidth float64

	// StrokeWidth is the line color (black by default)
	StrokeColor *Color

	// FillColor is color of the area which line creates (disabled by default)
	FillColor *Color

	// XFormatter is a custom value formatter for X values (values as is by default)
	XFormatter ValueFormatter

	// YFormatter is a custom value formatter for Y values (values as is by default)
	YFormatter ValueFormatter
}
