package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"bufio"
	"strings"
	"time"
)

func main () {
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

	time.Sleep(3 * time.Second)
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
		os.Exit(0)
	}
}
