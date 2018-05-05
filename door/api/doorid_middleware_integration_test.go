// +build integration

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/tlmiller/garage-door-controller/door"
	"github.com/tlmiller/garage-door-controller/door/dummy"
)

func TestDoorIsFound(t *testing.T) {
	middleware := NewDoorIdMiddleware(door.DoorServiceFunc(
		func(id door.Id) door.Door {
			if id != door.Id("test") {
				t.Errorf("middleware door id requesting wrong id %s", id)
			}
			return dummy.NewDoor(door.Id("test"))
		}))

	req := httptest.NewRequest("GET", "/api/door/test", nil)
	res := httptest.NewRecorder()

	urlVars := make(map[string]string)
	urlVars[doorIdKey] = "test"
	req = mux.SetURLVars(req, urlVars)

	middleware.Middleware(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
			rDoor, ok := req.Context().Value(doorIdKey).(door.Door)
			if ok == false {
				t.Error("failed converting middleware door to door.Door")
			}

			if rDoor.Id() != door.Id("test") {
				t.Error("Recieved door is not the one asked for")
			}
		})).ServeHTTP(res, req)
}

func TestDoorIdNotFound(t *testing.T) {
	middleware := NewDoorIdMiddleware(door.DoorServiceFunc(
		func(id door.Id) door.Door {
			if id != door.Id("noexist") {
				t.Errorf("middleware door id requesting wrong id %s", id)
			}
			return nil
		}))

	req := httptest.NewRequest("GET", "/api/door/noexist", nil)
	res := httptest.NewRecorder()

	urlVars := make(map[string]string)
	urlVars[doorIdKey] = "noexist"
	req = mux.SetURLVars(req, urlVars)

	middleware.Middleware(http.HandlerFunc(
		func(res http.ResponseWriter, req *http.Request) {
		})).ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("middleware door id not found did not return 404 instead %d", res.Code)
	}
}
