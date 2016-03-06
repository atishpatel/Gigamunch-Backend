package errors

import (
	"bytes"
	"fmt"
)

// TODO test in playground

// Errors is returned by operations when there are multiple possible errors for
// the function.
type Errors []error

func (m Errors) Error() string {
	numErrors := 0
	var buffer bytes.Buffer

	for _, e := range m {
		if e != nil {
			buffer.WriteString(e.Error() + ";")
			numErrors++
		}
	}
	return fmt.Sprintf("%d errors: %s", numErrors, buffer.String())
}

// HasErrors returns true if there is a none nil error in the Errors
func (m Errors) HasErrors() bool {
	for _, e := range m {
		if e != nil {
			return true
		}
	}
	return false
}

// AddError adds an error to Errors
func (m Errors) AddError(err error) {
	m = append(m, err)
}
