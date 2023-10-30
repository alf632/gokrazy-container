package main

import (
	"log"

	"github.com/alf632/gokrazy-container/podmanManager"
)

func main() {
	container := podmanManager.NewPodmanInstance(
		"victron",
		"localhost/gokrazy-victron",
		"latest",
		"/etc/localtime:/etc/localtime:ro,/run:/run,/opt/victronConnect/config.yml:/victron/config.yml:ro",
		true,
		true)
	container.AddBuildContext("/opt/victronConnect")
	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}
