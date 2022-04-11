package server

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	httpErrReasonKey  = "error_reason"
	httpErrMessageKey = "error_message"

	notFoundMsg            = "not found"
	invalidInputMsg        = "invalid input yaml"
	internalServerErrorMsg = "internal server error"

	errorInvalidSpec = "InvalidSpec"
)

// ValidationError is the interface for validation error
type ValidationError interface {
	Error() string
	HTTPCode() int32
	Reason() string
	Raw() error
}

type errImpl struct {
	err    error
	code   int32
	reason string
}

func (e *errImpl) Error() string {
	return e.err.Error()
}

func (e *errImpl) HTTPCode() int32 {
	return e.code
}

func (e *errImpl) Reason() string {
	return e.reason
}

func (e *errImpl) Raw() error {
	return e.err
}

func NewInvalidSpec(err error) ValidationError {
	return &errImpl{
		err:    err,
		code:   http.StatusBadRequest,
		reason: errorInvalidSpec,
	}
}

func handleValidationError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp[httpErrReasonKey] = invalidInputMsg
	resp[httpErrMessageKey] = fmt.Sprintf("%+v", err)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("failed to json marshal response")
	}
	w.Write(jsonResp)
}

func handleNotFoundError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp[httpErrReasonKey] = notFoundMsg
	resp[httpErrMessageKey] = fmt.Sprintf("%+v", err)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("failed to json marshal response")
	}
	w.Write(jsonResp)
}

func handleInternalError(w http.ResponseWriter, err error, errorMsg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp[httpErrReasonKey] = internalServerErrorMsg
	resp[httpErrMessageKey] = fmt.Sprintf("%s error: %+v", errorMsg, err)
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
