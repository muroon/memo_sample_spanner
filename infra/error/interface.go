package apperror

import "memo_sample_spanner/domain/app"

// ErrorManager error manager interface
type ErrorManager interface {
	Wrap(err error, code app.Error) error
	LogMessage(err error) string
	Code(err error) int
}
