package main

import "fmt"
import "github.com/stianeikeland/go-rpio"
import "os"

func main () {
	fmt.Println("Hello world!!")

	err := rpio.Open()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pin := rpio.Pin(10)

	pin.Output()
	pin.Low()
	
}

