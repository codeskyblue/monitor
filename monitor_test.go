package monitor

import (
	"testing"
)

func TestNCpu(t *testing.T) {
	if Ncpu() == 0 {
		t.Error("Ncpu cannot be 0, but result is 0")
	}
}
