package main

import (
	"os/exec"
	"os"
)

func main () {
	cmd := exec.Command("fluidsynth", "-a", "alsa", "-g", "6", "/usr/share/sounds/sf2/FluidR3_GM.sf2")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}
