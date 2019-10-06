package apperrorsub

import (
	"fmt"
	"memo_sample_spanner/infra/error"

	"github.com/srvc/fail"
)

// NewErrorManager new error Manager
func NewErrorManager() apperror.ErrorManager {
	return errorManager{}
}

type errorManager struct{}

func (em errorManager) Wrap(err error, code int) error {
	return fail.Wrap(
		err,
		fail.WithCode(code),
		fail.WithIgnorable(),
	)
}

func (em errorManager) LogMessage(err error) string {
	return fmt.Sprintf("%T\nCode:%d\nStackTrace:%+v\n",
		err,
		fail.Unwrap(err).Code,
		fail.Unwrap(err).StackTrace,
	)
}

func (em errorManager) Code(err error) int {
	return fail.Unwrap(err).Code.(int)
}
