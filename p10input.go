package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
	"io"
	"bufio"
)

func main () {
	fmt.Println("Hello world!!")

	rperr := rpio.Open()
	if rperr != nil {
		fmt.Println(rperr)
		os.Exit(1)
	}

	pin := rpio.Pin(10)

	pin.Input()
	current := pin.Read()
	wait := 2 * time.Millisecond
	go func () {
		for {
			time.Sleep(wait)
			now := pin.Read()
			if current != now {
				state := "High"
				if now == rpio.Low {
					state = "Low"
				}
				fmt.Println("Change :" + state)
				current = now
			}
		}		
	}()

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
}

