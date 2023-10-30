package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alf632/gokrazy-container/podmanManager"
	"golang.org/x/sys/unix"
)

func main() {
	initBluetooth()
	container := podmanManager.NewPodmanInstance(
		"bluetooth",
		"localhost/gokrazy-bluetooth",
		"latest",
		"/etc/localtime:/etc/localtime:ro,/run:/run,/opt/bluetooth/bt-agent.conf:/opt/bt-agent.conf",
		true,
		true)
	container.AddBuildContext("/opt/bluetooth")
	if err := container.Run(); err != nil {
		log.Fatal(err)
	}
}

func initBluetooth() error {
	// modprobe the hci_uart driver for Raspberry Pi (3B+, others)
	for _, mod := range []string{
		"kernel/crypto/ecc.ko",
		"kernel/crypto/ecdh_generic.ko",
		"kernel/net/bluetooth/bluetooth.ko",
		"kernel/drivers/bluetooth/btbcm.ko",
		"kernel/drivers/bluetooth/hci_uart.ko",
	} {
		if err := loadModule(mod); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	dev := "hci0"
	target, err := checkBluetoothInterface(dev)
	if err != nil {
		log.Printf("Bluetooth interface %v not found.", target)
	} else {
		fmt.Printf("Bluetooth device %v: %v\n", dev, target)
	}

	return nil
}

func checkBluetoothInterface(device string) (string, error) {
	target, err := os.Readlink("/sys/class/bluetooth/hci0")
	if err != nil {
		return "", fmt.Errorf("Bluetooth interface %v not found", device)
	}
	return target, nil
}

func loadModule(mod string) error {
	f, err := os.Open(filepath.Join("/lib/modules", release, mod))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := unix.FinitModule(int(f.Fd()), "", 0); err != nil {
		if err != unix.EEXIST &&
			err != unix.EBUSY &&
			err != unix.ENODEV &&
			err != unix.ENOENT {
			return fmt.Errorf("FinitModule(%v): %v", mod, err)
		}
	}
	modname := strings.TrimSuffix(filepath.Base(mod), ".ko")
	log.Printf("modprobe %v", modname)
	return nil
}

var release = func() string {
	var uts unix.Utsname
	if err := unix.Uname(&uts); err != nil {
		fmt.Fprintf(os.Stderr, "minitrd: %v\n", err)
		os.Exit(1)
	}
	return string(uts.Release[:bytes.IndexByte(uts.Release[:], 0)])
}()
