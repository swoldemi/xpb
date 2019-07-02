package browser

import (
	"errors"

	"github.com/tebeka/selenium"
)

var (
	// ErrInvalidReply is returned when an expected type assertion fails on an ExectueScript call
	ErrInvalidReply = errors.New("browser: ExecuteScript returned unexpected reply")
)

// ReadyStateCond will return true if the session's DOM ready state is complete.
// This function is used as a Selenium wait condition
func ReadyStateCond (wd selenium.WebDriver) (bool, error) {
	// Reference: https://github.com/tebeka/selenium/blob/master/remote_test.go#L1333
	script := "return document.readyState"
	reply, err := wd.ExecuteScript(script, []interface{}{})
	if err != nil {
		return false, err
	}
	result, ok := reply.(string)
	if !ok {
		return false, ErrInvalidReply
	}
	for result != "complete" {
		return false, nil

	}
	return true, nil
}
