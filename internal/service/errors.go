package service

import "errors"

var ErrNotFound = errors.New("not found")
var ErrPermissionDenied = errors.New("permission denied")
