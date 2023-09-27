package main

import (
	"log"

	"github.com/alf632/gokrazy-container/podmanManager"
)

func main() {
	container := podmanManager.NewPodmanInstance(
		"pigpiod",
		"zinen2/alpine-pigpiod",
		"latest",
		"/etc/localtime:/etc/localtime:ro",
		true,
		false)
	container.AddDevice("/dev/gpiochip0")
	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}
