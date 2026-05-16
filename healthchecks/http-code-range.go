package healthchecks

import (
	"fmt"
	"simple-docker-healthcheck/constants"
)

type HTTPCodeHealthCheckRangeStatusCode struct {
	URL string

	MinExpectedStatusCode int
	MaxExpectedStatusCode int
}

func (h *HTTPCodeHealthCheckRangeStatusCode) Dump() string {
	return fmt.Sprintf("HTTP code healthcheck on URL '%s' (expected status code range '%d-%d')", h.URL, h.MinExpectedStatusCode, h.MaxExpectedStatusCode)
}

func (h *HTTPCodeHealthCheckRangeStatusCode) Execute() (bool, error) {
	statusCode, _, err := retrieveHTTPResponse(h.URL)
	if err != nil {
		return false, err
	}

	if statusCode >= h.MinExpectedStatusCode && statusCode <= h.MaxExpectedStatusCode {
		return true, nil
	} else {
		return false, fmt.Errorf("expected status code between %d and %d, got %d", h.MinExpectedStatusCode, h.MaxExpectedStatusCode, statusCode)
	}
}

func (h *HTTPCodeHealthCheckRangeStatusCode) AreParametersValid() (bool, []string) {
	var errors []string
	if err := isURLValid(h.URL); err != nil {
		errors = append(errors, err.Error())
	}
	if h.MinExpectedStatusCode == constants.HTTP_STATUS_CODE_UNSET {
		errors = append(errors, "minimum expected status code is required")
	}
	if h.MaxExpectedStatusCode == constants.HTTP_STATUS_CODE_UNSET {
		errors = append(errors, "maximum expected status code is required")
	}
	if h.MinExpectedStatusCode != constants.HTTP_STATUS_CODE_UNSET && h.MaxExpectedStatusCode != constants.HTTP_STATUS_CODE_UNSET && h.MinExpectedStatusCode > h.MaxExpectedStatusCode {
		errors = append(errors, "minimum expected status code cannot be greater than maximum expected status code")
	}
	return len(errors) == 0, errors
}
