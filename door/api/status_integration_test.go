// +build integration

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tlmiller/garage-door-controller/door"
	"github.com/tlmiller/garage-door-controller/door/dummy"
)

func TestStatusHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/door/test", nil)
	res := httptest.NewRecorder()

	testDoor := dummy.NewDoor(door.Id("test"))
	req = req.WithContext(context.WithValue(req.Context(), doorIdKey, testDoor))

	doorStatusHandler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", res.Code)
	}

	if res.HeaderMap.Get("content-type") != "application/json" {
		t.Errorf("wrong content-type header of %s", res.HeaderMap.Get("content-type"))
	}
}
