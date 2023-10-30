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

type hostapdConf struct {
	SSID      string `json:"ssid"`
	PW        string `json:"passphrase"`
	Interface string `json:"interface"`
	Address   string `json:"address"`
}

func main() {
	container := podmanManager.NewPodmanInstance(
		"hostapd",
		"localhost/gokrazy-hostapd",
		"latest",
		true,
		true)
	container.AddVolume("/etc/localtime:/etc/localtime:ro", false)
	container.AddVolume("/run:/run", false)
	container.AddVolume("/perm/hostapd:/config", true)
	container.AddBuildContext("/opt/hostapd")

	configFile := flag.String("config", "/opt/hostapd/config.json", "path to config file")
	flag.Parse()
	var hConfig hostapdConf
	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatal("Error reading "+*configFile, err)
	}

	d := json.NewDecoder(strings.NewReader(string(data)))

	err = d.Decode(&hConfig)
	if err != nil {
		log.Fatal("Error parsing config", err)
	}

	//if _, err := os.Stat("/perm/hostapd/hostapd.conf"); errors.Is(err, os.ErrNotExist) {
	content := []byte(fmt.Sprintf(`
# "g" simply means 2.4GHz band
hw_mode=g
# the channel to use
channel=10
# limit the frequencies used to those allowed in the country
ieee80211d=1
# the country code
country_code=DE
# 802.11n support
ieee80211n=1
# QoS support, also required for full speed on 802.11n/ac/ax
wmm_enabled=1

# the name of the AP
ssid=%s
# 1=wpa, 2=wep, 3=both
auth_algs=1
# WPA2 only
wpa=2
wpa_key_mgmt=WPA-PSK
rsn_pairwise=CCMP
wpa_passphrase=%s`, hConfig.SSID, hConfig.PW))
	os.WriteFile("/perm/hostapd/hostapd.conf", content, 0644)
	//}

	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
	log.Fatal("container terminated")
}
