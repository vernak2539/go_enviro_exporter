package sensors

import (
	"github.com/rubiojr/go-enviroplus/mics6814"
)

// This is the MICS6814 sensor. Measure changes in the concentration of gases including carbon monoxide (CO),
// nitrogen dioxide (NO2), and ammonia (NH3) with this Breakout Garden compatible I2C breakout.

type MicsSensor struct {
	sensor *mics6814.Device
}

func NewMicsSensor() *MicsSensor {
	sensor, err := mics6814.New()
	if err != nil {
		panic(err)
	}

	go func() {
		sensor.StartReading()
	}()

	return &MicsSensor{
		sensor: sensor,
	}
}

func (s *MicsSensor) Close() {
	s.sensor.Halt()
}

func (s *MicsSensor) GetGasMeasurements() mics6814.Readings {
	return s.sensor.LastValue()
}