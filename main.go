package main

import (
	"fmt"
	"time"

	"github.com/rubiojr/go-enviroplus/ltr559"
)

func main() {
	d, err := ltr559.New()
	if err != nil {
		panic(err)
	}

	for {
		p, err := d.Proximity()
		if err != nil {
			panic(err)
		}

		fmt.Println("proximity: ", p)
		l, err := d.Lux()
		if err != nil {
			panic(err)
		}

		fmt.Println("      lux: ", l)
		time.Sleep(1 * time.Second)
	}

}
