package api

import (
	"log"
	"net/http"
)

type StatusError interface {
	error
	Status() int
}

type HandlerError struct {
	Err  error
	Code int
}

// implement error interface on StatusError
func (e HandlerError) Error() string {
	return e.Err.Error()
}

// returns http status code
func (e HandlerError) Status() int {
	return e.Code
}

type HttpErrorHandler func(http.ResponseWriter, *http.Request) error

func (handler HttpErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := handler(w, r)
	if err != nil {
		log.Print(err)
		switch e := err.(type) {
		case StatusError:
			http.Error(w, e.Error(), e.Status())
		default:
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
		}
	}
}
