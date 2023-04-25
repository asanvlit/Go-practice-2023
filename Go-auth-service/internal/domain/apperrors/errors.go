package apperrors

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrAlreadyRegisteredUserEmail = errors.New("user with the email already exists")
var ErrInvalidEmailFormat = errors.New("email validation failed")
var ErrInvalidPasswordFormat = errors.New("password validation failed")
var ErrInvalidRequestFormat = errors.New("invalid request format")
var ErrInvalidRequestBody = errors.New("invalid request body")
var ErrInvalidIdFormat = errors.New("invalid id format (not uuid)")
var ErrInvalidOffsetFormat = errors.New("invalid offset format")
var ErrInvalidLimitFormat = errors.New("invalid limit format")
var ErrInvalidDateFormat = errors.New("invalid date format")

var ErrInternalJsonProcessing = errors.New("failed to process json")
var ErrNatsPublishing = errors.New("failed to publish message to NATS")
var ErrDbQueryProcessing = errors.New("failed to execute query to db")
