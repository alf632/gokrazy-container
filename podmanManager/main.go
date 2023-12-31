package podmanManager

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gokrazy/gokrazy"
)

type PodmanInstance struct {
	name        string
	image       string
	tag         string
	hostNetwork bool
	privileged  bool
	volumes     []string
	envVars     []string

	devices []string

	buildContext string
}

func NewPodmanInstance(name, image, tag string, hostNetwork, privileged bool) *PodmanInstance {
	return &PodmanInstance{name: name, image: image, tag: tag, hostNetwork: hostNetwork, privileged: privileged}
}

func (pi *PodmanInstance) AddVolume(volume string, mkdir bool) {
	println("adding Volume: ", volume)
	if mkdir {
		dir := strings.Split(volume, ":")[0]
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Println("Creating Dir", dir)
			mkdirP(dir)
		}
	}
	pi.volumes = append(pi.volumes, volume)
}

func mkdirP(dir string) {
	if err := syscall.Mkdir(dir, 0777); err != nil {
		if err == syscall.ENOENT {
			volSplit := strings.Split(dir, "/")
			if len(volSplit) > 1 {
				parent := strings.Join(volSplit[:len(volSplit)-1], "/")
				log.Println("creating parent", parent)
				mkdirP(parent)
			} else {
				log.Println("giving up recursion")
			}
			return
		} else if err != syscall.EEXIST {
			log.Fatal(err)
		}
		// The directory already exists
		log.Printf("directory exists")
	}
}

func (pi *PodmanInstance) AddDevice(device string) {
	println("adding Device: ", device)
	pi.devices = append(pi.devices, device)
}

func (pi *PodmanInstance) AddBuildContext(bcontext string) {
	println("adding BuildContext:", bcontext)
	pi.buildContext = bcontext
}

func (pi *PodmanInstance) AddEnv(envVar string) {
	println("adding BuildContext:", envVar)
	pi.envVars = append(pi.envVars, envVar)
}

func (pi PodmanInstance) checkImageExists() bool {
	cmd := exec.Command("/usr/local/bin/podman", "images")
	cmd.Env = expandPath(os.Environ())
	cmd.Env = append(cmd.Env, "TMPDIR=/tmp")
	out, err := cmd.CombinedOutput()
	outStr := string(out)
	if err != nil {
		log.Fatalf("looking up images failed with %s\n", err)
	}
	exists := strings.Contains(outStr, pi.image)
	fmt.Println("image exists:", exists)
	return exists
}

func (pi PodmanInstance) build() error {
	log.Println("building image", pi.image)
	if err := podman("build", "--no-cache",
		"-t", pi.image+":"+pi.tag,
		pi.buildContext); err != nil {
		return err
	}
	return nil
}

func (pi PodmanInstance) Run() error {
	if err := mountVar(); err != nil {
		return err
	}

	exists := pi.checkImageExists()
	if !exists {
		// wait for ntp aka. internet connection
		gokrazy.WaitForClock()
		if len(pi.buildContext) > 0 {
			if err := pi.build(); err != nil {
				return err
			}
		}
	}

	if err := podman("kill", pi.name); err != nil {
		log.Print(err)
	}

	if err := podman("rm", pi.name); err != nil {
		log.Print(err)
	}

	startArgs := []string{"run", "-td"}
	for _, device := range pi.devices {
		startArgs = append(startArgs, "--device", device)
	}
	for _, volume := range pi.volumes {
		startArgs = append(startArgs, "-v", volume)
	}
	for _, envVar := range pi.envVars {
		startArgs = append(startArgs, "-e", envVar)
	}
	if pi.hostNetwork {
		startArgs = append(startArgs, "--network", "host")
	}
	if pi.privileged {
		startArgs = append(startArgs, "--privileged")
	}
	startArgs = append(startArgs, "--name", pi.name, pi.image+":"+pi.tag)

	log.Println("starting Container with:", startArgs)
	if err := podman(startArgs...); err != nil {
		return err
	}

	if err := podman("logs", "-f", pi.name); err != nil {
		return err
	}
	return nil
}

// podman wraps the podman binary and redirects STDIO
func podman(args ...string) error {
	podman := exec.Command("/usr/local/bin/podman", args...)
	podman.Env = expandPath(os.Environ())
	podman.Env = append(podman.Env, "TMPDIR=/tmp")
	podman.Stdin = os.Stdin
	podman.Stdout = os.Stdout
	podman.Stderr = os.Stderr
	if err := podman.Run(); err != nil {
		return fmt.Errorf("%v: %v", podman.Args, err)
	}
	return nil
}

// mountVar bind-mounts /perm/container-storage to /var if needed.
// This could be handled by an fstab(5) feature in gokrazy in the future.
func mountVar() error {
	b, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		log.Printf("Cannot Check mountpoint!")
		return err
	}
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		mountpoint := parts[4]
		log.Printf("Found mountpoint %q", parts[4])
		if mountpoint == "/var" {
			log.Printf("/var file system already mounted, nothing to do")
			return nil
		}
	}

	if err := syscall.Mkdir("/perm/container-storage", 0600); err != nil {
		if err != syscall.EEXIST {
			return err
		}
		// The directory already exists
		log.Printf("directory already exists")
	}

	if err := syscall.Mount("/perm/container-storage", "/var", "", syscall.MS_BIND, ""); err != nil {
		return fmt.Errorf("mounting /perm/container-storage to /var: %v", err)
	}

	return nil
}

// expandPath returns env, but with PATH= modified or added
// such that both /user and /usr/local/bin are included, which podman needs.
func expandPath(env []string) []string {
	extra := "/user:/usr/local/bin"
	found := false
	for idx, val := range env {
		parts := strings.Split(val, "=")
		if len(parts) < 2 {
			continue // malformed entry
		}
		key := parts[0]
		if key != "PATH" {
			continue
		}
		val := strings.Join(parts[1:], "=")
		env[idx] = fmt.Sprintf("%s=%s:%s", key, extra, val)
		found = true
	}
	if !found {
		const busyboxDefaultPATH = "/usr/local/sbin:/sbin:/usr/sbin:/usr/local/bin:/bin:/usr/bin"
		env = append(env, fmt.Sprintf("PATH=%s:%s", extra, busyboxDefaultPATH))
	}
	return env
}
