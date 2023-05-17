package compute

import (
	"log"
	"testing"

	"github.com/hallllll/miraiweed/ready"
)

func TestAllForOneSheet(t *testing.T) {
	t.Log("gogogog")
	log.Println("oioioi")
	paths_test, err := ready.NewPATHs()
	paths_test.ChangeBase("../data")
	paths_test.ChangeStorage("2023_05_17_074409_823")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(paths_test.StudentFolder())
	err = AllForOneSheet(paths_test)
	if err != nil {
		t.Fatal(err)
	}
}
