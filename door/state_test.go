package door

import (
	"testing"
)

type TestProvider struct {
	closed bool
	open   bool
}

func (p *TestProvider) IsClosed() bool {
	return p.closed
}

func (p *TestProvider) IsOpen() bool {
	return p.open
}

func TestConstructionState(t *testing.T) {
	tests := []struct {
		provider StateProvider
		state    State
	}{
		{&TestProvider{true, false}, DoorClosed},
		{&TestProvider{false, true}, DoorOpen},
		{&TestProvider{true, true}, DoorUnknown},
		{&TestProvider{false, false}, DoorUnknown},
	}

	for _, test := range tests {
		sMachine := NewStateMachine(test.provider)

		if current := sMachine.Current(); current != test.state {
			t.Errorf("new state machine state '%s != %s for open == %t & closed == %t",
				current, test.state, test.provider.IsOpen(), test.provider.IsClosed())
		}
		if toState, isTrans := sMachine.IsTransitioning(); isTrans != false && toState != DoorUnknown {
			t.Error("new statemachine was not isTransitioning() == false, nil")
		}
	}
}

func TestToggleStates(t *testing.T) {
	tests := []struct {
		// init provider state
		initClosed, initOpen bool
		// trans provider state
		transClosed, transOpen bool
		// trans state
		transFrom, transTo State
		// post state & error
		postState State
	}{
		// Closed -> Open
		{true, false, false, true, DoorClosed, DoorOpen, DoorOpen},
		// Open -> Closed
		{false, true, true, false, DoorOpen, DoorClosed, DoorClosed},
		// Unknown -> Open
		{false, false, false, true, DoorUnknown, DoorUnknown, DoorOpen},
		// Unknown -> Closed
		{false, false, true, false, DoorUnknown, DoorUnknown, DoorClosed},
		// Unknown -> Unknown
		{false, false, false, false, DoorUnknown, DoorUnknown, DoorUnknown},
	}

	for _, test := range tests {
		provider := &TestProvider{test.initClosed, test.initOpen}
		sMachine := NewStateMachine(provider)
		sync := make(chan struct{}, 0)

		err := sMachine.Toggle(func(from State, to State) {
			if from != test.transFrom {
				t.Errorf("from state mismatch %s != %s", from, test.transFrom)
			}
			if to != test.transTo {
				t.Errorf("to state mismatch %s != %s", to, test.transTo)
			}
			provider.closed = test.transClosed
			provider.open = test.transOpen
			sync <- struct{}{}
		})

		if err != nil {
			t.Errorf("unexpected error for transition: %s", err)
		}

		<-sync
		if sMachine.Current() != test.postState {
			t.Errorf("stat mismatch at then end of transitiong %s != %s", sMachine.Current(), test.postState)
		}
	}
}

func TestIsTransitioning(t *testing.T) {
	tests := []struct {
		provider StateProvider
	}{
		{&TestProvider{false, true}},
	}

	for _, test := range tests {
		sMachine := NewStateMachine(test.provider)
		sMachine.Toggle(func(from State, to State) {
			if toState, isTrans := sMachine.IsTransitioning(); isTrans == false || toState != to {
				t.Errorf("StateMachine.IsTransitioning() %t != true || %s != %s",
					isTrans, toState, to)
			}
		})
	}
}
