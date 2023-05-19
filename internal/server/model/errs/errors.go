package errs

import "errors"

var ErrUserAlreadyExist = errors.New("user already exist")
var ErrUserNotFound = errors.New("user not found")
var ErrResNotFound = errors.New("resource not found")

var ErrTokenNotFound = errors.New("unauthorized")
var ErrTokenInvalid = errors.New("invalid token")
