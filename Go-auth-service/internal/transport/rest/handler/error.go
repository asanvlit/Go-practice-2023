package handler

import (
	"Golang-practice-2023/internal/domain/apperrors"
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
	case apperrors.ErrUserNotFound:
		err := myHttp.WriteResponse(Error{Code: 404, Message: err.Error()}, w, contentType, http.StatusNotFound)
		return err
	case apperrors.ErrInvalidEmailFormat, apperrors.ErrInvalidPasswordFormat, apperrors.ErrInvalidRequestFormat,
		apperrors.ErrInvalidRequestBody, apperrors.ErrInvalidIdFormat, apperrors.ErrAlreadyRegisteredUserEmail:
		err := myHttp.WriteResponse(Error{Code: 400, Message: err.Error()}, w, contentType, http.StatusBadRequest)
		return err
	case apperrors.ErrInternalJsonProcessing, apperrors.ErrNatsPublishing, apperrors.ErrDbQueryProcessing:
		err := myHttp.WriteResponse(Error{Code: 500, Message: "Failed to execute"}, w, contentType, http.StatusInternalServerError)
		return err
	}
	return nil
}
