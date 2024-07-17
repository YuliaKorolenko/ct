package app

import "errors"

var ErrorNoRights = errors.New("no rights to change Ad")
var ErrorValidation = errors.New("wrong arguments")
var ErrorNoAd = errors.New("no such add, before create it")
var ErrorNoSuchUser = errors.New("this user don't exist")
var ErrorServer = errors.New("server error")

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

type NoSuchAdError struct {
	Err error
}

func (v NoSuchAdError) Error() string {
	return v.Err.Error()
}

type NoSuchUserError struct {
	Err error
}

func (v NoSuchUserError) Error() string {
	return v.Err.Error()
}
