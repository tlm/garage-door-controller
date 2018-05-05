package dummy

import (
	"time"

	"github.com/tlmiller/garage-door-controller/door"
)

type Door struct {
	DoorId        door.Id
	StateMachine  *door.StateMachine
	StateProvider StateProvider
}

type StateProvider struct {
	Closed bool
	Open   bool
}

func (d *Door) Id() door.Id {
	return d.DoorId
}

func (s *StateProvider) IsClosed() bool {
	return s.Closed
}

func (s *StateProvider) IsOpen() bool {
	return s.Open
}

func (d *Door) IsTriggered() (door.State, bool) {
	return d.StateMachine.IsTransitioning()
}

func NewDoor(id door.Id) *Door {
	nDoor := Door{
		DoorId:        id,
		StateMachine:  nil,
		StateProvider: StateProvider{true, false},
	}
	nDoor.StateMachine = door.NewStateMachine(&nDoor.StateProvider)
	return &nDoor
}

func (d *Door) State() door.State {
	return d.StateMachine.Current()
}

func (d *Door) Trigger() error {
	err := d.StateMachine.Toggle(func(from door.State, to door.State) {
		time.Sleep(time.Second * 1000)
		d.StateProvider.Closed = from == door.DoorOpen
		d.StateProvider.Open = from == door.DoorClosed
	})

	if err == door.ErrAlreadyTransitioning {
		return door.ErrAlreadyTriggered
	}
	return err
}
