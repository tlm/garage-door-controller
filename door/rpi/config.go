package rpi

import (
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

func Cleanup() {
	err := rpio.Close()
	if err != nil {
		panic(err)
	}
}

func Init() error {
	err := rpio.Open()
	if err != nil {
		return fmt.Errorf("failed to init rpi: %s", err)
	}
	return nil
}
