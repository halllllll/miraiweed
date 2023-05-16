package compute

import (
	"testing"

	"github.com/hallllll/miraiweed/ready"
)

func TestAllForOneSheet(t *testing.T) {
	paths_test, err := ready.NewPATHs()
	if err != nil {
		t.Fatal(err)
	}

	AllForOneSheet(paths_test)
}
