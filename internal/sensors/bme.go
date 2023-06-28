package sensors

import (
	"log"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

// This is the BME sensor, which is used to measure temperature, humidity, and pressure.

type BmeSensor struct {
	sensor *bmxx80.Dev
	bus    i2c.BusCloser
}

func NewBmeSensor() *BmeSensor {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}

	sensor, err := bmxx80.NewI2C(bus, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}

	return &BmeSensor{
		sensor: sensor,
		bus:    bus,
	}
}

func (s *BmeSensor) Close() {
	s.bus.Close()
	s.sensor.Halt()
}

// GetHumidity Return humidity converted to relative humidity
func (s *BmeSensor) GetHumidity() float64 {
	e := physic.Env{}

	if err := s.sensor.Sense(&e); err != nil {
		log.Fatal(err)
		return 0
	}

	return float64(e.Humidity) / float64(physic.PercentRH)
}

// GetPressure Return pressure converted from nanopascal to hectopascal
func (s *BmeSensor) GetPressure() float64 {
	e := physic.Env{}

	if err := s.sensor.Sense(&e); err != nil {
		log.Fatal(err)
		return 0
	}

	return float64(e.Pressure) / float64(physic.KiloPascal/10)
}

func (s *BmeSensor) GetTemperature() float64 {
	e := physic.Env{}

	if err := s.sensor.Sense(&e); err != nil {
		log.Fatal(err)
		return 0
	}

	return e.Temperature.Celsius()
}
