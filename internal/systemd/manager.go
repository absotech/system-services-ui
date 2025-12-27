package systemd

import (
	"errors"
	"os/exec"
	"strings"
)

func unit(name string) string {
	return name + ".service"
}

func Restart(name string) error {
	if !AllowedServices[name] {
		return errors.New("service not allowed")
	}

	cmd := exec.Command("/bin/systemctl", "restart", unit(name))
	return cmd.Run()
}

func Status(name string) (enabled bool, running bool, err error) {
	if !AllowedServices[name] {
		return false, false, errors.New("service not allowed")
	}

	enabledOut, err := exec.Command(
		"/bin/systemctl", "is-enabled", unit(name),
	).Output()

	if err == nil && strings.TrimSpace(string(enabledOut)) == "enabled" {
		enabled = true
	}

	activeOut, err := exec.Command(
		"/bin/systemctl", "is-active", unit(name),
	).Output()

	if err == nil && strings.TrimSpace(string(activeOut)) == "active" {
		running = true
	}

	return enabled, running, nil
}

func List() []string {
	services := make([]string, 0, len(AllowedServices))
	for s := range AllowedServices {
		services = append(services, s)
	}
	return services
}
