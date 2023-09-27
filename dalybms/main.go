package main

import (
	"log"

	"github.com/alf632/gokrazy-container/podmanManager"
)

func main() {
	container := podmanManager.NewPodmanInstance(
		"dalybms",
		"localhost/gokrazy-dalybms",
		"latest",
		"/etc/localtime:/etc/localtime:ro,/run:/run",
		true,
		true)
	container.AddBuildContext("/opt/dalybms")
	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}
