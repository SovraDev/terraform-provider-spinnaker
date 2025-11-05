package spinnaker

import (
	"fmt"
	"regexp"
)

func validateApplicationName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("Only alphanumeric characters or '-' allowed in %q", k))
	}
	return
}

func validatePipelineTriggerType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	validTypes := []string{
		"webhook",
	}
	isValid := false
	for _, t := range validTypes {
		if value == t {
			isValid = true
			break
		}
	}
	if !isValid {
		errors = append(errors, fmt.Errorf("Invalid value for %q: %q. Valid values are: %v", k, value, validTypes))
	}
	return
}

