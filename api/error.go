package api

import (
	"net/http"
)

func Error(res http.ResponseWriter, statusCode int) {
	res.WriteHeader(statusCode)
}
