package main

import (
	"log"

	"github.com/alf632/gokrazy-container/podmanManager"
)

func main() {
	container := podmanManager.NewPodmanInstance(
		"zigbee2mqtt",
		"koenkk/zigbee2mqtt",
		"latest",
		true,
		true)
	container.AddEnv("TZ=Europe/Berlin")
	container.AddVolume("/perm/zigbee2mqtt/data:/app/data", true)
	container.AddDevice("/dev/ttyUSB0")

	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}
