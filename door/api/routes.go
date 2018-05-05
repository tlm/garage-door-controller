package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tlmiller/garage-door-controller/door"
)

func GetRoutes(router *mux.Router) {
	sRouter := router.PathPrefix(fmt.Sprintf("/door/{%s}", doorIdKey)).Subrouter()

	sRouter.Use(NewDoorIdMiddleware(door.DefaultDirectory()).Middleware)
	sRouter.HandleFunc("", doorStatusHandler).Methods(http.MethodGet)
	sRouter.HandleFunc("/trigger", doorTriggerHandler).Methods(http.MethodPut)
}
