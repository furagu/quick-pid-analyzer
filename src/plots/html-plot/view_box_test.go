package htmlPlot

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestViewBoxWriter_Write(t *testing.T) {

	t.Run("should write viewBox after '<svg ' prefix is sent", func(t *testing.T) {
		w, buf := newWriter()

		_, err := w.Write([]byte("<svg foo bar"))
		assert.NoError(t, err)
		assert.Equal(t, `<svg viewBox="10 20 30 40" foo bar`, buf.String())
	})

	t.Run("should write viewBox when prefix is sent in smaller chunks", func(t *testing.T) {
		w, buf := newWriter()

		// first chunk
		_, err := w.Write([]byte("<sv"))
		assert.NoError(t, err)
		assert.Empty(t, buf.String(), "should keep the smaller content in cache, but not write it")

		// second chunk
		_, err = w.Write([]byte("g foo bar"))
		assert.NoError(t, err)
		assert.Equal(t, `<svg viewBox="10 20 30 40" foo bar`, buf.String(), "now it should write all the content with viewBox")
	})
}

func newWriter() (*viewBoxWriter, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return &viewBoxWriter{
		Target:  buf,
		ViewBox: [4]int{10, 20, 30, 40},
	}, buf
}
