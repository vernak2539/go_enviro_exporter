package main

import (
	"fmt"
	"log"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
	"periph.io/x/periph/host"
	"time"
)

func main() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Use i2creg I²C bus registry to find the first available I²C bus.
	b, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer b.Close()

	d, err := bmxx80.NewI2C(b, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}
	e := physic.Env{}

	for {
		if err := d.Sense(&e); err != nil {
			log.Fatal(err)
		}

		fmt.Print("    temp: ", int64(e.Temperature))
		fmt.Println()

		pressure := float64(e.Pressure) / float64(physic.KiloPascal/10)
		fmt.Print("pressure: ", pressure)
		fmt.Println()
		fmt.Print("pressure: ", int64(e.Pressure))
		fmt.Println()

		humidity := float64(e.Humidity) / float64(physic.PercentRH)
		fmt.Print("humidity: ", humidity)
		fmt.Println()
		fmt.Print("humidity: ", int64(e.Humidity))
		fmt.Println()

		time.Sleep(1 * time.Second)
	}
}
