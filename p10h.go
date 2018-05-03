package main

import "fmt"
import "github.com/stianeikeland/go-rpio"
import "os"

func main () {
	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pin := rpio.Pin(10)
	defer rpio.Close()

	pin.Output()
	pin.High()
	
}

