package systemd

import (
	"reflect"
	"testing"
)

func TestAllowlist(t *testing.T) {
	if err := SetAllowed([]string{"nginx", "budgetplanner.service", " transport-worker "}); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"nginx", "budgetplanner", "transport-worker"} {
		if !IsAllowed(name) {
			t.Errorf("expected %q to be allowed", name)
		}
	}
	if IsAllowed("sshd") {
		t.Error("sshd must not be allowed")
	}

	if !IsAllowedUnit("nginx.service") {
		t.Error("nginx.service should be allowed")
	}
	for _, u := range []string{"nginx", "nginx.socket", "sshd.service", "", ".service"} {
		if IsAllowedUnit(u) {
			t.Errorf("unit %q must not be allowed", u)
		}
	}

	want := []string{"budgetplanner", "nginx", "transport-worker"}
	if got := List(); !reflect.DeepEqual(got, want) {
		t.Errorf("List() = %v, want %v", got, want)
	}
}

func TestSetAllowedRejectsInvalid(t *testing.T) {
	for _, bad := range []string{"", "../etc", "a b", "x;rm"} {
		if err := SetAllowed([]string{bad}); err == nil {
			t.Errorf("expected error for %q", bad)
		}
	}
}
