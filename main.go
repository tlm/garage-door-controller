package main

import (
	"log"

	"github.com/tlmiller/garage-door-controller/door"
	//"github.com/github.com/tlmiller/garage-door-controller/door/dummy"
	"github.com/tlmiller/garage-door-controller/door/rpi"
	"github.com/tlmiller/garage-door-controller/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)
	err := rpi.Init()
	defer rpi.Cleanup()
	if err != nil {
		panic(err)
	}

	err = door.DefaultDirectory().Add(
		rpi.NewDoor(door.Id("static"), rpi.Pin(12), rpi.Pin(5), rpi.Pin(6)))
	//err := door.DefaultDirectory().Add(
	//	dummy.NewDoor(door.Id("test")))
	if err != nil {
		panic(err)
	}
	appServer := server.InitServer()
	appServer.Start()
}
