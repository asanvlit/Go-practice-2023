package handler

import (
	myErrors "Golang-practice-2023/internal/domain/apperrors"
	myHttp "Golang-practice-2023/pkg/utils/http"
	"github.com/pkg/errors"
	"net/http"
)

type Error struct {
	Code          int32
	Message       string
	DetailMessage string
}

func HandleError(w http.ResponseWriter, err error) error {
	switch errors.Cause(err) {
	case myErrors.ErrUserNotFound:
		return myHttp.BuildResponse(Error{Code: 404, Message: err.Error()}, w, contentType, http.StatusNotFound)
	}
	return nil
}
