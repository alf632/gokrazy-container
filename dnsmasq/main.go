package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alf632/gokrazy-container/podmanManager"
)

type dnsmasqConfig struct {
	IP        string `json:"ip"`
	IPRange   string `json:"ip-range"`
	DNS       string `json:"dns"`
	Domain    string `json:"domain"`
	Interface string `json:"interface"`
}

func main() {
	container := podmanManager.NewPodmanInstance(
		"dnsmasq",
		"localhost/gokrazy-dnsmasq",
		"latest",
		true,
		true)
	container.AddVolume("/etc/localtime:/etc/localtime:ro", false)
	container.AddVolume("/run:/run", false)
	container.AddVolume("/perm/dnsmasq:/etc/dnsmasq.d", true)
	container.AddBuildContext("/opt/dnsmasq")

	configFile := flag.String("config", "/opt/dnsmasq/config.json", "path to config file")
	flag.Parse()
	var dConfig dnsmasqConfig
	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatal("Error reading "+*configFile, err)
	}

	d := json.NewDecoder(strings.NewReader(string(data)))

	err = d.Decode(&dConfig)
	if err != nil {
		log.Fatal("Error parsing config", err)
	}

	//if _, err := os.Stat("/perm/dnsmasq/dnsmasq.conf"); errors.Is(err, os.ErrNotExist) {
	content := []byte(fmt.Sprintf(`
server=8.8.8.8
listen-address=127.0.0.1,%s
domain-needed
bogus-priv
filterwin2k
domain=%s
dhcp-range=%s
dhcp-option=option:dns-server,%s`, dConfig.IP, dConfig.Domain, dConfig.IPRange, dConfig.DNS))
	os.WriteFile("/perm/dnsmasq/dnsmasq.conf", content, 0644)
	//}

	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
	log.Fatal("container terminated")
}
