package entities

import (
	"fmt"
	"net/http"
	"runtime"
)

// HttpError represents an HTTP error with additional context.
type HttpError struct {
	Code    int            // HTTP status code
	Message any            // Error message (string or error)
	Caller  string         // Caller function and line number
	InLog   map[string]any // Additional log context
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

func newError(err any, code int, log ...map[string]any) HttpError {
	// NOTE: IF already an HttpError thus preserve its context.
	if he, ok := err.(HttpError); ok {
		if len(log) > 0 {
			for k, v := range log[0] { // Merge new log into existing one if any
				if he.InLog == nil {
					he.InLog = make(map[string]any)
				}
				he.InLog[k] = v
			}
		}
		return he
	}

	pc, _, line, _ := runtime.Caller(2)
	details := runtime.FuncForPC(pc)

	e := HttpError{
		Code:    code,
		Message: err,
		Caller:  fmt.Sprintf("%s#%d", details.Name(), line),
	}
	if len(log) > 0 {
		e.InLog = log[0]
	}
	return e
}

func ErrorBadRequest(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusBadRequest, log...)
}

func ErrorUnprocessableEntity(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusUnprocessableEntity, log...)
}

func ErrorUnauthorized(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusUnauthorized, log...)
}

func ErrorNotImplemented() HttpError {
	return newError("Not Implemented", http.StatusNotImplemented)
}

func ErrorForbidden(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusForbidden, log...)
}

func ErrorNotAcceptable(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusNotAcceptable, log...)
}

func ErrorNotFound(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusNotFound, log...)
}

func ErrorConflict(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusConflict, log...)
}

func ErrorPreconditionFailed(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusPreconditionFailed, log...)
}

func ErrorInternal(err any, log ...map[string]any) HttpError {
	return newError(err, http.StatusInternalServerError, log...)
}
