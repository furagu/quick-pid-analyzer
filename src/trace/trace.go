package trace

import (
	"fmt"
	"strconv"

	"github.com/furagu/quick-pid-analyzer/src/bblog"
	"github.com/furagu/quick-pid-analyzer/src/utils"
)

type Trace struct {
	Name  string
	PTerm float64
	Len   int

	Time []float64
	Gyro []float64
	P    []float64
}

func TracesFromLogSession(s *bblog.Session) (*Trace, *Trace, *Trace, error) {
	roll := &Trace{Name: "Roll", PTerm: float64(s.Header.RollP), Len: len(s.Values)}
	pitch := &Trace{Name: "Pitch", PTerm: float64(s.Header.PitchP), Len: len(s.Values)}
	yaw := &Trace{Name: "Yaw", PTerm: float64(s.Header.YawP), Len: len(s.Values)}

	dataPointsTotal := len(s.Values)
	for _, t := range []*Trace{roll, pitch, yaw} {
		t.Time = make([]float64, dataPointsTotal, dataPointsTotal)
		t.Gyro = make([]float64, dataPointsTotal, dataPointsTotal)
		t.P = make([]float64, dataPointsTotal, dataPointsTotal)
	}

	oTime := s.FieldOffset["time (us)"]

	var oRollGyro, oPitchGyro, oYawGyro int
	if _, ok := s.FieldOffset["gyroADC[0]"]; ok {
		oRollGyro = s.FieldOffset["gyroADC[0]"]
		oPitchGyro = s.FieldOffset["gyroADC[1]"]
		oYawGyro = s.FieldOffset["gyroADC[2]"]
	} else if _, ok := s.FieldOffset["gyroData[0]"]; ok {
		oRollGyro = s.FieldOffset["gyroData[0]"]
		oPitchGyro = s.FieldOffset["gyroData[1]"]
		oYawGyro = s.FieldOffset["gyroData[2]"]
	} else if _, ok := s.FieldOffset["ugyroADC[0]"]; ok {
		oRollGyro = s.FieldOffset["ugyroADC[0]"]
		oPitchGyro = s.FieldOffset["ugyroADC[1]"]
		oYawGyro = s.FieldOffset["ugyroADC[2]"]
	}

	oRollP := s.FieldOffset["axisP[0]"]
	oPitchP := s.FieldOffset["axisP[1]"]
	oYawP := s.FieldOffset["axisP[2]"]

	for i, v := range s.Values {
		time, err := strconv.ParseFloat(v[oTime], 64)
		if err != nil {
			return nil, nil, nil, err
		}
		time *= 1e-6
		roll.Time[i] = time
		pitch.Time[i] = time
		yaw.Time[i] = time

		rollGyro, err := strconv.ParseFloat(v[oRollGyro], 64)
		if err != nil {
			return nil, nil, nil, err
		}
		roll.Gyro[i] = rollGyro

		pitchGyro, err := strconv.ParseFloat(v[oPitchGyro], 64)
		if err != nil {
			return nil, nil, nil, err
		}
		pitch.Gyro[i] = pitchGyro

		yawGyro, err := strconv.ParseFloat(v[oYawGyro], 64)
		if err != nil {
			return nil, nil, nil, err
		}
		yaw.Gyro[i] = yawGyro

		rollP, err := strconv.ParseFloat(v[oRollP], 64)
		if err != nil {
			return nil, nil, nil, err
		}
		roll.P[i] = rollP

		pitchP, err := strconv.ParseFloat(v[oPitchP], 64)
		if err != nil {
			return nil, nil, nil, err
		}
		pitch.P[i] = pitchP

		yawP, err := strconv.ParseFloat(v[oYawP], 64)
		if err != nil {
			return nil, nil, nil, err
		}
		yaw.P[i] = yawP
	}

	return roll, pitch, yaw, nil
}

func (t *Trace) Print() {
	fmt.Printf("Trace: %s\nPTerm: %f\n", t.Name, t.PTerm)
	fmt.Printf("Time: %s\n", utils.FloatsToString(t.Time))
	fmt.Printf("Gyro: %s\n", utils.FloatsToString(t.Gyro))
	fmt.Printf("P: %s\n", utils.FloatsToString(t.P))
}
