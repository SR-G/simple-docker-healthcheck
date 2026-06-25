package healthchecks

import (
	"fmt"

	"github.com/SR-G/sul"
)

type FileHealthCheck struct {
	FileName string
}

func (h *FileHealthCheck) Dump() string {
	return fmt.Sprintf("file healthcheck, checking the file '%s' as being available on filename", h.FileName)
}

func (h *FileHealthCheck) Execute() (bool, error) {
	if sul.IsFileAvailable(h.FileName) {
		return true, nil
	} else {
		return false, fmt.Errorf("expected file '%s' not found on filesystem", h.FileName)
	}
}

func (h *FileHealthCheck) AreParametersValid() (bool, []string) {
	var errors []string
	if h.FileName == "" {
		errors = append(errors, "filename to monitor can't be empty")
	}
	return len(errors) == 0, errors
}
