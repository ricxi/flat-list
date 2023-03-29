package task

import "errors"

var ErrMissingField = errors.New("missing field is required")
var ErrTaskNotFound = errors.New("task not found")
