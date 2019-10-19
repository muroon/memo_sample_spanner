package apperror

import (
	"fmt"
	"memo_sample_spanner/domain/app"
	"strconv"

	"github.com/morikuni/failure"
)

// NewErrorManager new error Manager
func NewErrorManager() ErrorManager {
	return errorManager{}
}

type errorManager struct{}

func (em errorManager) Wrap(err error, code app.Error) error {

	cd := failure.StringCode(fmt.Sprintf("%d", code))
	err = failure.Wrap(err, failure.WithCode(cd))

	return err
}

func (em errorManager) LogMessage(err error) string {
	return fmt.Sprintf("%T\nCode:%d\n%+v\n",
		err,
		em.Code(err),
		err,
	)
}

func (em errorManager) Code(err error) int {
	var code int
	codeVal, ok := failure.CodeOf(err)
	if ok {
		code, _ = strconv.Atoi(codeVal.ErrorCode())
	}
	return code
}
