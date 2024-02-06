package storage

import (
	"errors"
	"net/http"
	"os"

	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry/remote/errcode"
)

func IsErrorCode(err error, code string) bool {
	var ec errcode.Error
	return errors.As(err, &ec) && ec.Code == code
}

func ErrCode(err error) int {
	var ec *errcode.ErrorResponse
	if errors.As(err, &ec) {
		return ec.StatusCode
	}
	return http.StatusInternalServerError
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, os.ErrNotExist) || ErrCode(err) == http.StatusNotFound || errors.Is(err, errdef.ErrNotFound)
}

func Error(w http.ResponseWriter, err error) {
	var ec *errcode.ErrorResponse
	switch {
	case errors.Is(err, os.ErrExist):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, os.ErrNotExist):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrInvalidArtifactType):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.As(err, &ec):
		if len(ec.Errors) < 1 {
			http.Error(w, err.Error(), ec.StatusCode)
			return
		}
		http.Error(w, ec.Errors[len(ec.Errors)-1].Message, ec.StatusCode)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
