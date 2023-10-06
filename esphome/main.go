package main

import (
	"log"

	"github.com/alf632/gokrazy-container/podmanManager"
)

func main() {
	container := podmanManager.NewPodmanInstance(
		"esphome",
		"ghcr.io/esphome/esphome",
		"latest",
		true,
		false)
	container.AddEnv("TZ=Europe/Berlin")
	container.AddVolume("/etc/localtime:/etc/localtime:ro", false)
	container.AddVolume("/perm/esphome:/config", true)

	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}
