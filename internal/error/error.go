package error

import "fmt"

// Error represents a error in the application
type Error struct {
	Layer  string `json:"layer,omitempty"`
	Reason string `json:"reason,omitempty"`
	Inner  error  `json:"inner,omitempty"`

	Code int `json:"-"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Inner != nil {
		return fmt.Sprintf("[%v error: (%s) -> %v]", e.Layer, e.Reason, e.Inner)
	}
	return fmt.Sprintf("[%v error: (%s)]", e.Layer, e.Reason)
}

func (e Error) Unwrap() error {
	return e.Inner
}

// Is implemets errors.Is() evaluator function
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return (e.Reason == t.Reason)
}

func (e Error) String() string {
	return e.Error()
}

// Clone clones the error
func (e Error) Clone() *Error {
	return e.CloneWithInner(e.Inner)
}

// CloneWithInner clones the error with provided inner error
func (e Error) CloneWithInner(err error) *Error {
	return New(e.Layer, e.Reason, err)
}

// IsMatchesCode checks wether the two errors have the same code
/* func (e Error) IsMatchesCode(err *Error) bool {
	return e.Code == err.Code
} */

// New returns a new application error
func New(layer, reason string, inner error) *Error {
	return NewWithCode(layer, reason, 0, inner)
}

// NewWithCode returns a new application error with a code
func NewWithCode(layer, reason string, code int, inner error) *Error {
	return &Error{
		Layer:  layer,
		Reason: reason,
		Inner:  inner,
		Code:   code,
	}
}
