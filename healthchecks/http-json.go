package healthchecks

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

type HTTPHealthCheckJSONPath struct {
	URL string

	JSONPath      string
	ExpectedValue string
	Insensitive   bool
}

func (h *HTTPHealthCheckJSONPath) Dump() string {
	return fmt.Sprintf("HTTP JSON healthcheck on URL '%s' (expected value '%s' at JSONPath '%s', insensitive: '%t')", h.URL, h.ExpectedValue, h.JSONPath, h.Insensitive)
}

func (h *HTTPHealthCheckJSONPath) Execute() (bool, error) {
	_, response, err := retrieveHTTPResponse(h.URL)
	if err != nil {
		return false, err
	}

	jsonObject := interface{}(nil)
	err = json.Unmarshal([]byte(response), &jsonObject)
	if err != nil {
		return false, fmt.Errorf("response is not a valid JSON: %w", err)
	}

	result, err := jsonpath.Get(h.JSONPath, jsonObject)
	if err != nil {
		return false, err
	}

	responseValue, ok := result.(string)
	if !ok {
		return false, fmt.Errorf("found content value at JSONPath %s is not a string %s", h.JSONPath, string(responseValue))
	}

	if h.Insensitive {
		responseValue = strings.ToLower(responseValue)
		h.ExpectedValue = strings.ToLower(h.ExpectedValue)
	}

	if responseValue == h.ExpectedValue {
		return true, nil
	} else {
		return false, fmt.Errorf("expected value '%s' at JSONPath %s, got '%s'", h.ExpectedValue, h.JSONPath, responseValue)
	}
}

func (h *HTTPHealthCheckJSONPath) AreParametersValid() (bool, []string) {
	var errors []string
	if err := isURLValid(h.URL); err != nil {
		errors = append(errors, err.Error())
	}
	if h.JSONPath == "" {
		errors = append(errors, "JSONPath is required")
	}
	if h.ExpectedValue == "" {
		errors = append(errors, "expected value is required")
	}
	return len(errors) == 0, errors
}
