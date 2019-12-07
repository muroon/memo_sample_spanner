package app

// Error error code
type Error int

const (
	// InputError user input error
	InputError Error = iota + 1
	// DBError db error
	DBError
)

// ContextKey key for transaction context
type ContextKey string

const (
	// ResKey response key
	ResKey = "http.response.key"

	// SpannerTransaction spanner transaction key
	SpannerTransaction = "spanner.transaction"
)
