package door

import (
	"fmt"
)

type Id string

var (
	ErrAlreadyTriggered = fmt.Errorf("door is already triggered")
)

type Door interface {
	Id() Id
	IsTriggered() (State, bool)
	State() State
	Trigger() error
}
