package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorStatus(t *testing.T) {
	tests := []struct {
		status int
	}{
		{http.StatusBadRequest},
		{http.StatusNotFound},
		{http.StatusInternalServerError},
		{http.StatusNotImplemented},
	}

	for _, test := range tests {
		res := httptest.NewRecorder()
		Error(res, test.status)
		if res.Code != test.status {
			t.Errorf("error status mismatch %d != %d", res.Code, test.status)
		}
	}
}
