package errors

import "fmt"

// ErrNotImplemented is returned when a requested operation is not implemented
var ErrNotImplemented error = fmt.Errorf("not implemented")
