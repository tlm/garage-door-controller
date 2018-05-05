package rpi

import (
	"time"

	"github.com/stianeikeland/go-rpio"

	"github.com/tlmiller/garage-door-controller/door"
)

type Door struct {
	id         door.Id
	sMachine   *door.StateMachine
	triggerPin rpio.Pin
}

type Pin uint8

type stateProvider struct {
	closedPin rpio.Pin
	openPin   rpio.Pin
}

func NewDoor(id door.Id, trigPin Pin, closedPin Pin, openPin Pin) *Door {
	door := Door{
		id:         id,
		sMachine:   door.NewStateMachine(newStateProvider(closedPin, openPin)),
		triggerPin: rpio.Pin(trigPin),
	}
	door.triggerPin.Output()
	door.triggerPin.Low()
	return &door
}

func newStateProvider(closedPin Pin, openPin Pin) door.StateProvider {
	sProvider := stateProvider{
		closedPin: rpio.Pin(closedPin),
		openPin:   rpio.Pin(openPin),
	}
	sProvider.closedPin.Input()
	sProvider.closedPin.PullUp()
	sProvider.openPin.Input()
	sProvider.openPin.PullUp()
	return &sProvider
}

func (d *Door) Id() door.Id {
	return d.id
}

func (s *stateProvider) IsClosed() bool {
	return s.closedPin.Read() == rpio.Low
}

func (s *stateProvider) IsOpen() bool {
	return s.openPin.Read() == rpio.Low
}

func (d *Door) IsTriggered() (door.State, bool) {
	return d.sMachine.IsTransitioning()
}

func (d *Door) State() door.State {
	return d.sMachine.Current()
}

func (d *Door) Trigger() error {
	err := d.sMachine.Toggle(func(from door.State, to door.State) {
		d.triggerPin.High()
		time.Sleep(time.Second * 1)
		d.triggerPin.Low()
	})

	if err == door.ErrAlreadyTransitioning {
		return door.ErrAlreadyTriggered
	}
	return err
}
