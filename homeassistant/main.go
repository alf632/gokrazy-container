package main

import (
	"log"

	"github.com/alf632/gokrazy-container/podmanManager"
)

func main() {
	container := podmanManager.NewPodmanInstance(
		"homeassistant",
		"ghcr.io/home-assistant/home-assistant",
		"stable",
		true,
		false)
	container.AddEnv("TZ=Europe/Berlin")
	container.AddVolume("/etc/localtime:/etc/localtime:ro", false)
	container.AddVolume("/run:/run", false)
	container.AddVolume("/perm/ha:/config", true)

	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}
