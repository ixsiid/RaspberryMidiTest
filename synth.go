package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"bufio"
	"strings"
	"time"
	"github.com/stianeikeland/go-rpio"
)

type PinListener struct {
	High string
	Low  string
	Pin rpio.Pin
	Led rpio.Pin
}

type Pins [] PinListener

func main () {
	rperr := rpio.Open()
	if rperr != nil {
		fmt.Println(rperr)
		os.Exit(1)
	}

	var selectors Pins
	selectors = append(selectors, PinListener {
		Pin : rpio.Pin(8),
		High: "select 0 1 0 0",
		Low : "",
		Led : rpio.Pin(12),
	})
	selectors = append(selectors, PinListener {
		Pin : rpio.Pin(9),
		High: "select 0 1 0 48",
		Low : "",
		Led : rpio.Pin(12),
	})
	selectors = append(selectors, PinListener {
		Pin : rpio.Pin(10),
		High: "select 0 2 0 0",
		Low : "",
		Led : rpio.Pin(12),
	})
	selectors = append(selectors, PinListener {
		Pin : rpio.Pin(11),
		High: "select 0 2 0 16",
		Low : "",
		Led : rpio.Pin(12),
	})

	cmd := exec.Command("fluidsynth", "-a", "alsa", "-g", "6", "/usr/share/sounds/sf2/FluidR3_GM.sf2")
	stdin, _ := cmd.StdinPipe()
	cmdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr
	cmd.Start()

	fmt.Println("Waiting fluidsynt")
	cmdScanner := bufio.NewScanner(cmdout)
	for {
		if cmdScanner.Scan() {
			m := cmdScanner.Text()
			fmt.Println(m)
			if strings.Index(m, "help") >= 0 {
				fmt.Println("found 'help' word")
				break
			}
		} else {
			fmt.Println("Close scanner")
			break
		}
	}

	time.Sleep(1500 * time.Millisecond)
	fmt.Println("Started fluidsynth")

	go func () {
		for cmdScanner.Scan() {
			m := cmdScanner.Text()
			fmt.Println(m)
			if m == "cheers!" {
				fmt.Println("byebye")
				os.Exit(0)
			}
		}
	}()

// load sound font
	fonts := [] string { "/usr/share/sounds/sf2/TimGM6mb.sf2" }
	for _, f := range fonts {
		io.WriteString(stdin, "load " + f + "\n")
	}


	wait := 5 * time.Millisecond
	for _, item := range selectors {
		go func (selector PinListener) {
			fmt.Printf(">>> Pin %s is initializing\n", selector)
			selector.Pin.Input()
			selector.Led.PullDown()
			current := selector.Pin.Read()
			for {
				time.Sleep(wait)
				now := selector.Pin.Read()
				if current != now {
					current = now
					if current == rpio.High {
						fmt.Println(selector.High)
						io.WriteString(stdin, selector.High + "\n")
					} else {
						io.WriteString(stdin, selector.Low + "\n")
					}
				}
			}
		}(item)
	}

	var err error
	reader := bufio.NewReaderSize(os.Stdin, 4096)

	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		fmt.Print("input> ")
		fmt.Println(line)

		
		io.WriteString(stdin, line)
	}
	if err != io.EOF {
		panic(err)
	} else {
		fmt.Println("\nbyebye!")
	}

	for _, p := range selectors {
		p.Pin.PullDown()
	}
	rpio.Close()
	os.Exit(0)
}
