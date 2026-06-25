package healthchecks

import (
	"fmt"
	"strings"

	"github.com/SR-G/sul"
)

type FileContentHealthCheck struct {
	FileName        string
	ExpectedContent string
	Insensitive     bool
}

func (h *FileContentHealthCheck) Dump() string {
	return fmt.Sprintf("file content healthcheck, checking the file '%s' having the content '%s', similar to a grep (insensitive : %v)", h.FileName, h.ExpectedContent, h.Insensitive)
}

func (h *FileContentHealthCheck) Execute() (bool, error) {
	if !sul.IsFileAvailable(h.FileName) {
		return false, fmt.Errorf("expected file '%s' not found on filesystem", h.FileName)
	} else {
		content, err := sul.ReadContent(h.FileName)
		if err != nil {
			return false, fmt.Errorf("can't read content of file '%s' : %v", h.FileName, err)
		}
		s := strings.TrimSpace(string(content[:]))
		if s == "" {
			return false, fmt.Errorf("expected file '%s' is empty", h.FileName)
		}
		if h.Insensitive {
			if sul.ContainsI(s, h.ExpectedContent) {
				return true, nil
			} else {
				return false, fmt.Errorf("expected content '%s' not found in file '%s' (insensitive = true)", h.ExpectedContent, h.FileName)
			}
		} else {
			if strings.Contains(s, h.ExpectedContent) {
				return true, nil
			} else {
				return false, fmt.Errorf("expected content '%s' not found in file '%s' (insensitive = false)", h.ExpectedContent, h.FileName)
			}
		}
	}
}

func (h *FileContentHealthCheck) AreParametersValid() (bool, []string) {
	var errors []string
	if h.FileName == "" {
		errors = append(errors, "filename of file to monitor can't be empty")
	}
	if h.ExpectedContent == "" {
		errors = append(errors, "content to be grepped can't be empty")
	}
	return len(errors) == 0, errors
}
