package htmlPlot

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

var writeViewBoxAfter = []byte("<svg ")

// viewBoxWriter is a proxy writer with intercepts piped bytes and adds "viewBox"
// attribute to an <svg> tag.
// It receives bytes and caches them until the moment when string "<svg " has been
// received, then it adds the viewBox and pipes the rest of content as is.
// The viewBox attribute is needed to make SVG elements responsive in browser.
type viewBoxWriter struct {
	Target      io.Writer
	ViewBox     [4]int
	viewBoxSent bool
	sentBytes   []byte
}

// Write writes SVG content and appends "viewBox" attribute to found SVG tag
func (w *viewBoxWriter) Write(p []byte) (n int, err error) {
	// just proxy to the original writer if we already sent the viewBox attribute
	if w.viewBoxSent {
		return w.Target.Write(p)
	}

	// check if we already received the whole svg tag prefix
	// if not, then just cache the content and continue
	prefixLen := len(writeViewBoxAfter)
	newBytes := append(w.sentBytes, p...)
	if len(newBytes) < prefixLen {
		w.sentBytes = newBytes
		return 0, nil
	}

	// check if the received content is actually svg tag
	if !w.startsWith(newBytes, writeViewBoxAfter) {
		return 0, errors.Errorf("expected content to start with '%s' to write viewBox SVG attribute, actual value: '%s'", writeViewBoxAfter, newBytes)
	}

	// write part of the received content which ends after the tag name
	n, err = w.Target.Write(newBytes[0:prefixLen])
	if err != nil {
		return n, err
	}

	// append the viewBox attribute
	viewBox := fmt.Sprintf(`viewBox="%d %d %d %d" `, w.ViewBox[0], w.ViewBox[1], w.ViewBox[2], w.ViewBox[3])
	n, err = w.Target.Write([]byte(viewBox))
	if err != nil {
		return n, err
	}
	w.viewBoxSent = true

	// write the rest of received content
	return w.Target.Write(newBytes[prefixLen:])
}

// startsWith checks if the slice start with subSlice
func (w *viewBoxWriter) startsWith(slice []byte, subSlice []byte) bool {
	for key, value := range subSlice {
		if slice[key] != value {
			return false
		}
	}
	return true
}
