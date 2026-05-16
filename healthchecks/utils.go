package healthchecks

import (
	"fmt"
	"io"
	"net/http"
	"simple-docker-healthcheck/constants"
	"strings"
)

func isURLValid(url string) error {
	if url == "" {
		return fmt.Errorf("URL is required")
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("invalid URL '%s': must start with 'http://' or 'https://'", url)
	}
	return nil
}

func retrieveHTTPResponse(url string) (int, string, error) {

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

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
