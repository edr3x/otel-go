package entities

import (
	"fmt"
	"net/http"
	"runtime"
)

type HttpError struct {
	Code    int
	Message any
	Caller  string
}

func (se HttpError) Error() string {
	switch e := se.Message.(type) {
	case string:
		return e
	case error:
		return e.Error()
	default:
		return fmt.Sprintf("%v", se.Message)
	}
}

func (se HttpError) Status() int {
	return se.Code
}

func newError(err any, code int) HttpError {
	// NOTE: IF already an HttpError thus preserve its context.
	if he, ok := err.(HttpError); ok {
		return he
	}

	pc, _, line, _ := runtime.Caller(2)
	details := runtime.FuncForPC(pc)

	return HttpError{
		Code:    code,
		Message: err,
		Caller:  fmt.Sprintf("%s#%d", details.Name(), line),
	}
}

func ErrorBadRequest(err any) HttpError {
	return newError(err, http.StatusBadRequest)
}

func ErrorUnprocessableEntity(err any) HttpError {
	return newError(err, http.StatusUnprocessableEntity)
}

func ErrorUnauthorized(err any) HttpError {
	return newError(err, http.StatusUnauthorized)
}

func ErrorNotImplemented() HttpError {
	return newError("Not Implemented", http.StatusNotImplemented)
}

func ErrorForbidden(err any) HttpError {
	return newError(err, http.StatusForbidden)
}

func ErrorNotAcceptable(err any) HttpError {
	return newError(err, http.StatusNotAcceptable)
}

func ErrorNotFound(err any) HttpError {
	return newError(err, http.StatusNotFound)
}

func ErrorConflict(err any) HttpError {
	return newError(err, http.StatusConflict)
}

func ErrorPreconditionFailed(err any) HttpError {
	return newError(err, http.StatusPreconditionFailed)
}

func ErrorInternal(err any) HttpError {
	return newError(err, http.StatusInternalServerError)
}
