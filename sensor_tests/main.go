package main

import (
	"fmt"
	"github.com/rubiojr/go-enviroplus/pms5003"
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

	particulateMatterSensor, err := pms5003.New()
	if err != nil {
		panic(err)
	}
	go func() {
		particulateMatterSensor.StartReading()
	}()

	for {
		pm := particulateMatterSensor.LastValue()
		if err := d.Sense(&e); err != nil {
			log.Fatal(err)
		}

		fmt.Print("    temp: ", e.Temperature.Celsius())
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

		fmt.Println("-------")
		fmt.Println("PM1.0 ug/m3 (ultrafine):                        ", pm.Pm10Std)
		fmt.Println("PM2.5 ug/m3 (combustion, organic comp, metals): ", pm.Pm25Std)
		fmt.Println("PM10 ug/m3 (dust, pollen, mould spores):        ", pm.Pm100Std)
		fmt.Println("PM1.0 ug/m3 (atmos env):                        ", pm.Pm10Env)
		fmt.Println("PM2.5 ug/m3 (atmos env):                        ", pm.Pm25Env)
		fmt.Println("PM10 ug/m3 (atmos env):                         ", pm.Pm100Env)
		fmt.Println("0.3um 1 0.1L air:                               ", pm.Particles3um)
		fmt.Println("0.5um 1 0.1L air:                               ", pm.Particles5um)
		fmt.Println("1.0um 1 0.1L air:                               ", pm.Particles10um)
		fmt.Println("2.5um 1 0.1L air:                               ", pm.Particles25um)
		fmt.Println("5um 1 0.1L air:                                 ", pm.Particles50um)
		fmt.Println("10um 1 0.1L air:                                ", pm.Particles100um)

		time.Sleep(1 * time.Second)
	}
}
