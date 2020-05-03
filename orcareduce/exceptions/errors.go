package exceptions

import "fmt"

var (
	ErrActorsDepleted = fmt.Errorf("actors depleted")
)
