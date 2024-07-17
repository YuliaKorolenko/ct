package app

import "errors"

var ErrorNoRights = errors.New("no rights to change Ad")
var ErrorValidation = errors.New("wrong arguments")

type PermissionDeniedError struct {
	Err error
}

func (p PermissionDeniedError) Error() string {
	return p.Err.Error()
}

type ValidationApError struct {
	Err error
}

func (v ValidationApError) Error() string {
	return v.Err.Error()
}
