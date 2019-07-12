package util

import (
	"github.com/tebeka/selenium"
)


// ReadyStateCond will return true if the session's DOM ready state is complete.
// This function is used as a Selenium wait condition.
// Reference: https://github.com/tebeka/selenium/blob/master/remote_test.go#L1333
func ReadyStateCond(wd selenium.WebDriver) (bool, error) {
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


// ReadKeyFile reads a service account keyfile. Note:
// this application expects that the service account
// is activated before use. 
func ReadKeyFile(filepath string) map[string]interface{} {
	// Read key file only to extract project_id and client_email
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	var svcAccount map[string]interface{}
	err = json.Unmarshal(data, &svcAccount)
	if err != nil {
		log.Fatal(err)
	}

	return svcAccount
}
