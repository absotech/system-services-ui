package systemd

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

const unitSuffix = ".service"

var validName = regexp.MustCompile(`^[A-Za-z0-9:_.@-]+$`)

var (
	mu              sync.RWMutex
	allowedServices = map[string]bool{}
)

// NormalizeName validates a service name from the config and strips an
// optional ".service" suffix.
func NormalizeName(name string) (string, error) {
	n := strings.TrimSuffix(strings.TrimSpace(name), unitSuffix)
	if n == "" || !validName.MatchString(n) {
		return "", fmt.Errorf("invalid service name %q", name)
	}
	return n, nil
}

// SetAllowed replaces the allowlist with the given service names.
func SetAllowed(names []string) error {
	m := make(map[string]bool, len(names))
	for _, name := range names {
		n, err := NormalizeName(name)
		if err != nil {
			return err
		}
		m[n] = true
	}

	mu.Lock()
	allowedServices = m
	mu.Unlock()
	return nil
}

// IsAllowed reports whether the service name (without ".service") is allowed.
func IsAllowed(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return allowedServices[name]
}

// IsAllowedUnit reports whether the full unit name (e.g. "nginx.service") is allowed.
func IsAllowedUnit(unitName string) bool {
	if !strings.HasSuffix(unitName, unitSuffix) {
		return false
	}
	return IsAllowed(strings.TrimSuffix(unitName, unitSuffix))
}

func allowedList() []string {
	mu.RLock()
	defer mu.RUnlock()
	services := make([]string, 0, len(allowedServices))
	for s := range allowedServices {
		services = append(services, s)
	}
	sort.Strings(services)
	return services
}
