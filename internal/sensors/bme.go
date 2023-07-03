package sensors

import (
	"log"

	"github.com/rubiojr/go-enviroplus/bme280"
	"periph.io/x/conn/v3/physic"
)

// This is the BME sensor, which is used to measure temperature, humidity, and pressure.

type BmeSensor struct {
	sensor *bme280.BME280
}

func NewBmeSensor() *BmeSensor {
	sensor, err := bme280.New()
	if err != nil {
		log.Fatal(err)
	}

	return &BmeSensor{
		sensor: sensor,
	}
}

// GetHumidity Return humidity converted to relative humidity
func (s *BmeSensor) GetHumidity() float64 {
	e, _ := s.sensor.Read()

	return float64(e.Humidity) / float64(physic.PercentRH)
}

// GetPressure Return pressure converted from nanopascal to hectopascal
func (s *BmeSensor) GetPressure() float64 {
	e, _ := s.sensor.Read()

	return float64(e.Pressure) / float64(physic.KiloPascal/10)
}

func (s *BmeSensor) GetTemperature() float64 {
	e, _ := s.sensor.Read()

	return e.Temperature.Celsius()
}
