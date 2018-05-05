package door

import (
	"errors"
)

type Directory map[Id]Door

type DoorService interface {
	Find(id Id) Door
}

type DoorServiceFunc func(id Id) Door

func NewDirectory() Directory {
	return Directory{}
}

func (d Directory) Add(door Door) error {
	if d[door.Id()] != nil {
		return errors.New("door id already exists in directory")
	}
	d[door.Id()] = door
	return nil
}

func (d Directory) Find(id Id) Door {
	return d[id]
}

func (f DoorServiceFunc) Find(id Id) Door {
	return f(id)
}
