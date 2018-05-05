package door

import (
	"log"
)

type empty struct{}
type semaphore chan empty

type State string

type StateMachine struct {
	provider        StateProvider
	transitioning   bool
	transitioningTo State
	transitionLock  semaphore
}

type StateMachineError struct {
	err              string
	wasTransitioning bool
}

type StateProvider interface {
	IsClosed() bool
	IsOpen() bool
}

type Transition func(from State, to State)

const (
	DoorOpen    = "open"
	DoorClosed  = "closed"
	DoorUnknown = "unknown"
)

var (
	ErrAlreadyTransitioning = StateMachineError{
		err:              "already transitioning",
		wasTransitioning: true,
	}
)

func (s *StateMachine) Current() State {
	return providerState(s.provider)
}

func (e StateMachineError) Error() string {
	return e.err
}

func (s *StateMachine) IsTransitioning() (State, bool) {
	s.transitionLock <- empty{}
	if s.transitioning {
		return s.transitioningTo, true
	}
	<-s.transitionLock
	return s.nextState(), false
}

func (s *StateMachine) nextState() State {
	switch s.Current() {
	case DoorOpen:
		return DoorClosed
	case DoorClosed:
		return DoorOpen
	case DoorUnknown:
		return DoorUnknown
	default:
		panic("recieved an unknown state type")
	}
}

func NewStateMachine(provider StateProvider) *StateMachine {
	return &StateMachine{
		provider:       provider,
		transitioning:  false,
		transitionLock: make(semaphore, 1),
	}
}

func providerState(provider StateProvider) State {
	if provider.IsOpen() && provider.IsClosed() == false {
		return DoorOpen
	} else if provider.IsOpen() == false && provider.IsClosed() {
		return DoorClosed
	}
	return DoorUnknown
}

func (s *StateMachine) Toggle(trans Transition) error {
	s.transitionLock <- empty{}
	if s.transitioning {
		<-s.transitionLock
		return ErrAlreadyTransitioning
	}
	s.transitioning = true
	s.transitioningTo = s.nextState()
	<-s.transitionLock

	go func() {
		trans(s.Current(), s.transitioningTo)
		pState := providerState(s.provider)
		if s.transitioningTo != pState && s.transitioningTo != DoorUnknown {
			log.Printf("warning - transitioning resulted in mismatch state %s != %s",
				s.transitioningTo, pState)
		}
		s.transitionLock <- empty{}
		s.transitioning = false
		<-s.transitionLock
	}()

	return nil
}

func (e *StateMachineError) WasTransitioning() bool {
	return e.wasTransitioning
}
