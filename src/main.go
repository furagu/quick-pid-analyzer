package main

import (
	"log"
	"fmt"
	"context"
	"os"
	"github.com/furagu/quick-pid-analyzer/src/autodiscovery"
	"github.com/maxlaverse/blackbox-library/src/blackbox"
	"github.com/maxlaverse/blackbox-library/src/blackbox/stream"
)

func main() {
	newFiles := autodiscovery.GetFileChannel()

	readerOpts := blackbox.FlightLogReaderOpts{Raw: true}
	flightLog := blackbox.NewFlightLogReader(readerOpts)

	for filePath := range newFiles {
		log.Println("New file: " + filePath)
		readFile(filePath, flightLog)
	}
}

func readFile(path string, flightLog *blackbox.FlightLogReader) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	logFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	frameChan, err := flightLog.LoadFile(logFile, ctx)
	if err != nil {
		panic(err)
	}
	err = readFrames(frameChan)
	if err != nil {
		cancel()
		panic(err)
	}
}

func readFrames(frameChan <-chan blackbox.Frame) (err error) {
	for frame := range frameChan {
		// handle frame error
		err := frame.Error()
		if err != nil {
			if isErrorRecoverable(err) {
				log.Printf(`Frame '%s' with values %v has error: "%s"`, string(frame.Type()), frame.Values(), err.Error())
			} else {
				return err
			}
		}

		fmt.Println(frame)
	}
	return
}

func isErrorRecoverable(err error) bool {
	switch err.(type) {
	case *stream.ReadError, stream.ReadError:
		return false
	default:
		return true
	}
}
