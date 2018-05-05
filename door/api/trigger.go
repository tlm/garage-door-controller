package api

import (
	"encoding/json"
	"net/http"

	"github.com/tlmiller/garage-door-controller/door"
)

type TriggerResponse struct {
	StatusResponse
	err string `json:"error,omitempty"`
}

func doorTriggerHandler(res http.ResponseWriter, req *http.Request) {
	reqDoor := req.Context().Value(doorIdKey).(door.Door)
	triggerRes := TriggerResponse{
		StatusResponse: StatusResponse{
			Id:      reqDoor.Id(),
			Current: reqDoor.State(),
		},
	}

	err := reqDoor.Trigger()
	nextState, triggered := reqDoor.IsTriggered()
	triggerRes.IsTriggered = triggered
	triggerRes.Next = nextState

	res.Header().Set("content-type", "application/json")
	if err == door.ErrAlreadyTriggered {
		triggerRes.err = err.Error()
		res.WriteHeader(http.StatusConflict)
	} else if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		res.WriteHeader(http.StatusAccepted)
	}
	json.NewEncoder(res).Encode(triggerRes)
}
