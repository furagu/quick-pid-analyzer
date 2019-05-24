package main

import (
	"log"
	"context"
	"os"
	"github.com/furagu/quick-pid-analyzer/src/autodiscovery"
	"github.com/maxlaverse/blackbox-library/src/blackbox"
	"github.com/maxlaverse/blackbox-library/src/blackbox/stream"
)

type DataPoint struct {
	time int32
	gyroADC int32
	axisP0 int32
	axisP1 int32
	axisP2 int32
}


func main() {
	newFiles := autodiscovery.GetFileChannel()

	readerOpts := blackbox.FlightLogReaderOpts{Raw: true}
	flightLog := blackbox.NewFlightLogReader(readerOpts)

	for filePath := range newFiles {
		log.Println("Reading new file: " + filePath)
		dataPoints, err := readFile(filePath, flightLog)
		if err != nil {
			panic(err)
		}
		log.Printf("Found %d useful frames", len(dataPoints))
	}
}

func readFile(path string, flightLog *blackbox.FlightLogReader) (dataPoints []DataPoint, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	logFile, err := os.Open(path)
	defer logFile.Close()
	if err != nil {
		return
	}
	frameChan, err := flightLog.LoadFile(logFile, ctx)
	if err != nil {
		return
	}
	dataPoints, err = readFrames(frameChan, flightLog.FrameDef)
	if err != nil {
		cancel()
		return
	}
	return
}

func readFrames(frameChan <-chan blackbox.Frame, frameDef blackbox.LogDefinition) (dataPoints []DataPoint, err error) {
	for frame := range frameChan {
		// handle frame error
		err = frame.Error()
		if err != nil {
			if isErrorRecoverable(err) {
				log.Printf(`Frame '%s' with values %v has error: "%s"`, string(frame.Type()), frame.Values(), err.Error())
			} else {
				return
			}
		} else {
			switch frame.(type) {
			case *blackbox.MainFrame:
			  dataPoints = append(dataPoints, getDataPoint(frame, frameDef))
			default:
			}
		}
	}
	return
}

func getDataPoint(frame blackbox.Frame, frameDef blackbox.LogDefinition) (point DataPoint) {
	for k, v := range frame.Values().([]int32) {
		if i, _ := frameDef.GetFieldIndex(blackbox.FieldTime); k == i {
			point.time = v
			continue
		}
		if i, _ := frameDef.GetFieldIndex("axisP[0]"); k == i {
			point.axisP0 = v
			continue
		}
		if i, _ := frameDef.GetFieldIndex("axisP[1]"); k == i {
			point.axisP1 = v
			continue
		}
		if i, _ := frameDef.GetFieldIndex("axisP[2]"); k == i {
			point.axisP2 = v
			continue
		}
		if i, _ := frameDef.GetFieldIndex("gyroADC"); k == i {
			point.gyroADC= v
			continue
		}
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
