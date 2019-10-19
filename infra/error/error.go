package apperror

import (
	"fmt"
	"memo_sample_spanner/domain/app"

	"golang.org/x/xerrors"
)

// NewErrorManager new error Manager
func NewErrorManager() ErrorManager {
	return errorManager{}
}

type errorManager struct{}

type iWithCodeError interface {
	Code() int
}

type withCodeError struct {
	err   error
	code  app.Error
	frame xerrors.Frame
}

func (c withCodeError) Code() int {
	return int(c.code)
}

func (c withCodeError) Error() string {
	return fmt.Sprintf("code:%d, %#+v", c.code, c.err)
}

func (c withCodeError) Unwrap() error {
	return c.err
}

func (c withCodeError) FormatError(p xerrors.Printer) error {
	p.Print(c.Error())
	c.frame.Format(p)
	return nil
}
func (c withCodeError) Format(f fmt.State, ru rune) {
	xerrors.FormatError(c, f, ru)
}

func (em errorManager) Wrap(err error, code app.Error) error {
	e := xerrors.Errorf("error occurred: %w", err)
	return &withCodeError{err: e, code: code, frame: xerrors.Caller(1)}
}

func (em errorManager) LogMessage(err error) string {
	return err.Error()
}

func (em errorManager) Code(err error) int {
	var code int
	if e, ok := err.(iWithCodeError); ok {
		code = e.Code()
	}

	return code
}
