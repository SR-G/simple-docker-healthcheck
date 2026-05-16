package healthchecks

import (
	"fmt"
	"simple-docker-healthcheck/constants"
)

type HTTPCodeHealthCheckExactStatusCode struct {
	URL                string
	ExpectedStatusCode int
}

func (h *HTTPCodeHealthCheckExactStatusCode) Dump() string {
	return fmt.Sprintf("HTTP code healthcheck on URL '%s' (expected status code '%d')", h.URL, h.ExpectedStatusCode)
}

func (h *HTTPCodeHealthCheckExactStatusCode) Execute() (bool, error) {
	statusCode, _, err := retrieveHTTPResponse(h.URL)
	if err != nil {
		return false, err
	}

	if statusCode == h.ExpectedStatusCode {
		return true, nil
	} else {
		return false, fmt.Errorf("expected status code %d, got %d", h.ExpectedStatusCode, statusCode)
	}
}

func (h *HTTPCodeHealthCheckExactStatusCode) AreParametersValid() (bool, []string) {
	var errors []string
	if err := isURLValid(h.URL); err != nil {
		errors = append(errors, err.Error())
	}
	if h.ExpectedStatusCode == constants.HTTP_STATUS_CODE_UNSET {
		errors = append(errors, "expected status code is required")
	}
	return len(errors) == 0, errors
}
