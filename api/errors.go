package api

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// OAuthError is the JSON handler for OAuth2 error responses
type OAuthError struct {
	Err             string `json:"error"`
	Description     string `json:"error_description,omitempty"`
	InternalError   error  `json:"-"`
	InternalMessage string `json:"-"`
}

func (e *OAuthError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%s: %s", e.Err, e.Description)
}

// WithInternalError adds internal error information to the error
func (e *OAuthError) WithInternalError(err error) *OAuthError {
	e.InternalError = err
	return e
}

// WithInternalMessage adds internal message information to the error
func (e *OAuthError) WithInternalMessage(fmtString string, args ...interface{}) *OAuthError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

// Cause returns the root cause error
func (e *OAuthError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

// HTTPError is an error with a message and an HTTP status code.
type HTTPError struct {
	Code            int    `json:"code"`
	Message         string `json:"msg"`
	InternalError   error  `json:"-"`
	InternalMessage string `json:"-"`
	ErrorID         string `json:"error_id,omitempty"`
}

func (e *HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// Cause returns the root cause error
func (e *HTTPError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

// WithInternalError adds internal error information to the error
func (e *HTTPError) WithInternalError(err error) *HTTPError {
	e.InternalError = err
	return e
}

// WithInternalMessage adds internal message information to the error
func (e *HTTPError) WithInternalMessage(fmtString string, args ...interface{}) *HTTPError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

// ErrorCause provides error information
type ErrorCause interface {
	Cause() error
}

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	errorID := "0"
	switch e := err.(type) {
	case *HTTPError:
		if e.Code >= http.StatusInternalServerError {
			e.ErrorID = errorID
			// this will get us the stack trace too
		}
		if jsonErr := sendJSON(w, e.Code, e); jsonErr != nil {
			handleError(jsonErr, w, r)
		}
	case *OAuthError:
		if jsonErr := sendJSON(w, http.StatusBadRequest, e); jsonErr != nil {
			handleError(jsonErr, w, r)
		}
	case ErrorCause:
		handleError(e.Cause(), w, r)
	default:
		// hide real error details from response to prevent info leaks
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Infof("Encountered an unhandled error while servicing request. %+v\n", err)
		if _, writeErr := w.Write(
			[]byte(`{"code":500,"msg":"Internal server error","error_id":"` + errorID + `"}`),
		); writeErr != nil {
			logrus.Errorf("Encountered an unhandled error while writing response. %+v\n", writeErr)
		}
	}
}
