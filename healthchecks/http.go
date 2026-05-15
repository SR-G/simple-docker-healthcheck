package healthchecks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"simple-docker-healthcheck/constants"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

type HTTPCodeHealthCheckExactStatusCode struct {
	URL                string
	ExpectedStatusCode int
}

type HTTPCodeHealthCheckRangeStatusCode struct {
	URL string

	MinExpectedStatusCode int
	MaxExpectedStatusCode int
}

type HTTPHealthCheckText struct {
	URL string

	ExpectedText string
}

type HTTPHealthCheckJSONPath struct {
	URL string

	JSONPath      string
	ExpectedValue string
}

func retrieveHTTPResponse(url string) (int, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	req.Header.Set("User-Agent", constants.DEFAULT_USER_AGENT)

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(body), nil
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

func (h *HTTPHealthCheckText) Execute() (bool, error) {
	_, response, err := retrieveHTTPResponse(h.URL)
	if err != nil {
		return false, err
	}

	if strings.Contains(response, h.ExpectedText) {
		return true, nil
	} else {
		return false, fmt.Errorf("expected text %s not found in response", h.ExpectedText)
	}
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
	if responseValue == h.ExpectedValue {
		return true, nil
	} else {
		return false, fmt.Errorf("expected value '%s' at JSONPath %s, got '%s'", h.ExpectedValue, h.JSONPath, responseValue)
	}
}
