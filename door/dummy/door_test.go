package dummy

import (
	"testing"

	"github.com/tlmiller/garage-door-controller/door"
)

func TestDoorConstruction(t *testing.T) {
	dummyDoor := NewDoor(door.Id("static"))
	if dummyDoor.Id() != door.Id("static") {
		t.Error("door.Id() != \"static\"")
	}
}
