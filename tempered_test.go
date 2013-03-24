package tempered_test

import (
	"github.com/janne/tempered"
	"testing"
)

func TestNew(t *testing.T) {
	temp, err := tempered.New()

	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if len(temp.Devices) == 0 {
		t.Error("Devices is empty")
	}
}
