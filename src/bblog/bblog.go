package bblog

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

type Header struct {
	Firmware    string // fwType
	MaxThrottle int    // maxThrottle

	RollP int
	RollI int
	RollD int

	PitchP int
	PitchI int
	PitchD int

	YawP int
	YawI int
	YawD int
}

type Session struct {
	Header      *Header
	FieldOffset map[string]int
	Values      [][]string
}

const sampleDataPath = "sample.csv"

func SampleSession() (*Session, error) {
	h := &Header{
		Firmware:    "cleanflight",
		MaxThrottle: 2000,

		RollP: 42,
		RollI: 60,
		RollD: 35,

		PitchP: 46,
		PitchI: 70,
		PitchD: 38,

		YawP: 35,
		YawI: 100,
		YawD: 0,
	}

	fieldsToLoad := []string{
		"time (us)",
		"rcCommand[3]",
		"axisP[0]", "axisP[1]", "axisP[2]",
		"gyroADC[0]", "gyroADC[1]", "gyroADC[2]",
		"gyroData[0]", "gyroData[1]", "gyroData[2]",
		"ugyroADC[0]", "ugyroADC[1]", "ugyroADC[2]",
	}

	isRequiredField := map[string]bool{}
	for _, field := range fieldsToLoad {
		isRequiredField[field] = true
	}

	f, err := os.Open(sampleDataPath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	var inputFieldOffset map[string]int
	var outputFieldOffset map[string]int
	values := [][]string{}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if inputFieldOffset == nil {
			inputFieldOffset = make(map[string]int)

			for i, field := range record {
				if isRequiredField[field] {
					inputFieldOffset[field] = i
				}
			}

			if len(inputFieldOffset) == 0 {
				return nil, errors.New("Required fields not found in the CSV file")
			}

			outputFieldOffset = make(map[string]int)
			offset := 0
			for _, field := range fieldsToLoad {
				if _, ok := inputFieldOffset[field]; ok {
					outputFieldOffset[field] = offset
					offset++
				}
			}
		} else {
			v := make([]string, len(inputFieldOffset), len(inputFieldOffset))

			for field, i := range inputFieldOffset {
				j := outputFieldOffset[field]
				v[j] = record[i]
			}

			values = append(values, v)
		}

	}

	s := &Session{
		Header:      h,
		FieldOffset: outputFieldOffset,
		Values:      values,
	}

	return s, nil
}
