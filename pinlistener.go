package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
	"io"
	"bufio"
)

type PinListener struct {
	Pin  int
	High string
	Low  string
}

type Pins [] PinListener

func main () {
	fmt.Println("Hello world!!")

	rperr := rpio.Open()
	if rperr != nil {
		fmt.Println(rperr)
		os.Exit(1)
	}

	var pins Pins
	pins = append(pins, PinListener {
		Pin : 10,
		High: "pin 10 is change to high",
		Low : "pin 10 is change to low",
	})
	pins = append(pins, PinListener {
		Pin : 9,
		High: "higheeeeeerrrrrrrr",
		Low : "lo----------wwwww",
	})

	wait := 2 * time.Millisecond
	for _, item := range pins {
		go func (p PinListener) {
			fmt.Printf("Pin %d is initializing\n", p.Pin)
			pin := rpio.Pin(p.Pin)
			pin.Input()
			current := pin.Read()
			for {
				time.Sleep(wait)
				now := pin.Read()
				if current != now {
					current = now
					if current == rpio.High {
						fmt.Println(p.High)
					} else {
						fmt.Println(p.Low)
					}
				}
			}
		}(item)
	}

	var err error
	reader := bufio.NewReaderSize(os.Stdin, 4096)

	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		fmt.Println(line)
	}

	if err != io.EOF {
		panic(err)
	} else {
		fmt.Println("\nbyebye!")
		os.Exit(0)
	}

	for _, p := range pins {
		pin := rpio.Pin(p.Pin)
		pin.PullDown()
	}
}

