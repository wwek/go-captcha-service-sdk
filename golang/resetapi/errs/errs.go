package errs

import (
	"errors"
	"net/http"
)

// RemoteServiceErr ..
var RemoteServiceErr = errors.New("remote service error")

// MethodNotAllowedErr ..
var MethodNotAllowedErr = errors.New("method not allowed")

// NotFoundErr ..
var NotFoundErr = errors.New("status not found")

// BadRequestErr ..
var BadRequestErr = errors.New("bad request")

// RequestErr ..
var RequestErr = errors.New("request error")

// ParseDataErr ..
var ParseDataErr = errors.New("parse data error")

// UnauthorizedErr ..
var UnauthorizedErr = errors.New("unauthorized")

// ForbiddenErr ..
var ForbiddenErr = errors.New("forbidden")

// CheckBizCodeSuccess data.code == 200
func CheckBizCodeSuccess(code int64) bool {
	return code == http.StatusOK
}

// CheckBizCodeErr ..
func CheckBizCodeErr(code int64) error {
	if code == http.StatusBadRequest {
		return BadRequestErr
	} else if code == http.StatusMethodNotAllowed {
		return MethodNotAllowedErr
	} else if code == http.StatusNotFound {
		return NotFoundErr
	}

	return nil
}
