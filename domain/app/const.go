package app

// Error error code
type Error int

const (
	// InputError user input error
	InputError Error = iota + 1
	// DBError db error
	DBError
)
