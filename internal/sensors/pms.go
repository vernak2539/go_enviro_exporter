package sensors

import "github.com/rubiojr/go-enviroplus/pms5003"

// This is the PMS sensor, which is used to measure particulate matter.

type PmsSensor struct {
	sensor *pms5003.Device
}

func NewPmsSensor() *PmsSensor {
	sensor, err := pms5003.New()
	if err != nil {
		panic(err)
	}

	go func() {
		sensor.StartReading()
	}()

	return &PmsSensor{
		sensor: sensor,
	}
}

func (s *PmsSensor) GetPmMeasurement() *pms5003.PMS5003 {
	return s.sensor.LastValue()
}
