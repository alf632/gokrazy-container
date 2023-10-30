package main

import (
	"log"

	"github.com/alf632/gokrazy-container/podmanManager"
)

func main() {
	container := podmanManager.NewPodmanInstance(
		"iwd",
		"localhost/gokrazy-iwd",
		"latest",
		true,
		true)
	container.AddBuildContext("/opt/iwd")
	container.AddVolume("/etc/localtime:/etc/localtime:ro", false)
	container.AddVolume("/perm/resolv.conf:/etc/resolv.conf", false)
	container.AddVolume("/perm/iwd:/var/lib/iwd", true)

	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}
