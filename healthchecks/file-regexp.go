package healthchecks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/SR-G/sul"
)

type FileRegExpHealthCheck struct {
	FileName string
	RegExp   string
	r        *regexp.Regexp
}

func (h *FileRegExpHealthCheck) Dump() string {
	return fmt.Sprintf("file regexp healthcheck, checking the file '%s' matching the regexp '%s'", h.FileName, h.RegExp)
}

func (h *FileRegExpHealthCheck) Execute() (bool, error) {
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
		if h.r.MatchString(s) {
			return true, nil
		} else {
			return false, fmt.Errorf("file '%s' is not matching expected regexp '%s', actual content '%s'", h.FileName, h.RegExp, s)
		}
	}
}

func (h *FileRegExpHealthCheck) AreParametersValid() (bool, []string) {
	var errors []string
	if h.FileName == "" {
		errors = append(errors, "filename of file to monitor can't be empty")
	}
	if h.RegExp == "" {
		errors = append(errors, "regexp to be looked for in the file content can't be empty")
	}
	r, err := regexp.Compile(h.RegExp)
	if err != nil {
		errors = append(errors, "regexp '' can't be compiled : %v", h.RegExp, err.Error())
	}
	h.r = r
	return len(errors) == 0, errors
}
