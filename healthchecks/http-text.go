package healthchecks

import (
	"fmt"
	"strings"
)

type HTTPHealthCheckText struct {
	URL string

	ExpectedText string
	Insensitive  bool
}

func (h *HTTPHealthCheckText) Dump() string {
	return fmt.Sprintf("HTTP text healthcheck on URL '%s' (expected text '%s', insensitive: '%t')", h.URL, h.ExpectedText, h.Insensitive)
}

func (h *HTTPHealthCheckText) Execute() (bool, error) {
	_, response, err := retrieveHTTPResponse(h.URL)
	if err != nil {
		return false, err
	}

	if h.Insensitive {
		response = strings.ToLower(response)
		h.ExpectedText = strings.ToLower(h.ExpectedText)
	}

	if strings.Contains(response, h.ExpectedText) {
		return true, nil
	} else {
		return false, fmt.Errorf("expected text %s not found in response", h.ExpectedText)
	}
}

func (h *HTTPHealthCheckText) AreParametersValid() (bool, []string) {
	var errors []string
	if err := isURLValid(h.URL); err != nil {
		errors = append(errors, err.Error())
	}
	if h.ExpectedText == "" {
		errors = append(errors, "expected text is required")
	}
	return len(errors) == 0, errors
}
