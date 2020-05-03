package exceptions

import "fmt"

var (
	ErrActorsDepleted = fmt.Errorf("actors depleted")
	ErrPreconditionsNotMet = fmt.Errorf("preconditions not met")
)


type PreconditionError string

func (err PreconditionError) Error() string {
	return fmt.Sprintf("preconditions not met: %s", err)
}