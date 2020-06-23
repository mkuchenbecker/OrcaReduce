package exceptions

import (
	"fmt"
)

// PreconditionError is a runtime error that indicates all the preconditions to run a function are not yet satisfied.
// It should generally be used and regarded as a non-serious error that is expected during the normal operation
// of a service or state machine. For example, while its possible to guarentee all preconditions for a function are
// met at the time the function is called it's not possible to guarentee they will still be met at the time of
// preconditions are checked.
type PreconditionError string

// Error implements the error interface function so PreconditionError can be used
// as an error object.
func (err PreconditionError) Error() string {
	return fmt.Sprintf("preconditions not met: %s", string(err))
}
