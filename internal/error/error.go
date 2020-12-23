package error

import "fmt"

// Error represents a error in the application
type Error struct {
	Layer  string `json:"layer,omitempty"`
	Reason string `json:"reason,omitempty"`
	Code   int    `json:"code,omitempty"`
	Inner  error  `json:"inner,omitempty"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.Inner != nil {
		return fmt.Sprintf("[%v error: (%v:%s) -> %v]", e.Layer, e.Code, e.Reason, e.Inner)
	}
	return fmt.Sprintf("[%v error: (%v:%s)]", e.Layer, e.Code, e.Reason)
}

func (e Error) Unwrap() error {
	return e.Inner
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
	return New(e.Layer, e.Reason, e.Code, err)
}

// IsMatchesCode checks wether the two errors have the same code
func (e Error) IsMatchesCode(err *Error) bool {
	return e.Code == err.Code
}

// New returns a new application error
func New(layer, reason string, code int, inner error) *Error {
	return &Error{
		Layer:  layer,
		Reason: reason,
		Code:   code,
		Inner:  inner,
	}
}
