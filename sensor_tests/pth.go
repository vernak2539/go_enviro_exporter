package sensor_tests

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
	if err := d.Sense(&e); err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Printf("    temp: %8s", e.Temperature)
		fmt.Printf("pressure: %10s", e.Pressure)
		fmt.Printf("humidity: %9s", e.Humidity)
		time.Sleep(1 * time.Second)
	}
}
