package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tlmiller/garage-door-controller/api"
	"github.com/tlmiller/garage-door-controller/door"
)

type StatusResponse struct {
	Id          door.Id    `json:"id"`
	Current     door.State `json:"current"`
	IsTriggered bool       `json:"isTriggered"`
	Next        door.State `json:"next"`
}

func doorStatusHandler(res http.ResponseWriter, req *http.Request) {
	statDoor, ok := req.Context().Value(doorIdKey).(door.Door)

	if ok == false {
		log.Printf("failed converting context \"%s\" to door.Door", doorIdKey)
		api.Error(res, http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	nextState, triggered := statDoor.IsTriggered()
	status := StatusResponse{
		Id:          statDoor.Id(),
		Current:     statDoor.State(),
		IsTriggered: triggered,
		Next:        nextState,
	}
	json.NewEncoder(res).Encode(status)
}
