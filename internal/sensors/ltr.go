package sensors

import "github.com/rubiojr/go-enviroplus/ltr559"

// This is the LTR559 sensor, which is used to measure ambient light and proximity.

type LtrSensor struct {
	sensor *ltr559.LTR559
}

func NewLtrSensor() *LtrSensor {
	sensor, err := ltr559.New()
	if err != nil {
		panic(err)
	}

	return &LtrSensor{
		sensor: sensor,
	}
}

func (s *LtrSensor) GetProximity() float64 {
	proximity, err := s.sensor.Proximity()
	if err != nil {
		panic(err)
	}

	return proximity
}

func (s *LtrSensor) GetLux() float64 {
	lux, err := s.sensor.Lux()
	if err != nil {
		panic(err)
	}

	return lux
}
