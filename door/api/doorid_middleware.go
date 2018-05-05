package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tlmiller/garage-door-controller/api"
	"github.com/tlmiller/garage-door-controller/door"
)

type DoorIdMiddleware struct {
	service door.DoorService
}

func (d *DoorIdMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		doorId := door.Id(mux.Vars(req)[doorIdKey])
		door := d.service.Find(doorId)

		if door == nil {
			api.Error(res, http.StatusNotFound)
			return
		}
		newRequest := req.WithContext(context.WithValue(req.Context(), doorIdKey, door))
		next.ServeHTTP(res, newRequest)
	})
}

func NewDoorIdMiddleware(service door.DoorService) *DoorIdMiddleware {
	return &DoorIdMiddleware{
		service: service,
	}
}
