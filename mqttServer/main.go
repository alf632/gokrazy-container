package main

import (
	"log"
	"os"

	"github.com/alf632/gokrazy-container/podmanManager"
	"github.com/plus3it/gorecurcopy"
)

func main() {
	container := podmanManager.NewPodmanInstance(
		"mqttServer",
		"eclipse-mosquitto",
		"latest",
		true,
		false)
	container.AddEnv("TZ=Europe/Berlin")
	container.AddVolume("/etc/localtime:/etc/localtime:ro", false)
	container.AddVolume("/perm/mqtt/config:/mosquitto/config", true)
	container.AddVolume("/perm/mqtt/data:/mosquitto/data", true)
	container.AddVolume("/perm/mqtt/log:/mosquitto/log", true)

	if _, err := os.Stat("/perm/mqtt/config/mosquitto.conf"); os.IsNotExist(err) {
		if err := gorecurcopy.CopyDirectory("/opt/mosquitto", "/perm/mqtt/config/"); err != nil {
			log.Fatal(err)
		}
	}

	// touch logfile
	if file, err := os.OpenFile("/perm/mqtt/log/mosquitto.log", os.O_CREATE, 0666); err != nil {
		log.Fatal(err)
	} else {
		file.Close()
	}

	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}
