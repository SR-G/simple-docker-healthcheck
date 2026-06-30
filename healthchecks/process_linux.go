package healthchecks

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SR-G/sul"
)

func isProcessRunning(processName string, insensitive bool) (bool, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false, err
	}

	selfCmdlineBytes, err := os.ReadFile("/proc/self/cmdline")
	if err != nil {
		return false, err
	}
	selfCmdlineRaw := strings.TrimRight(string(selfCmdlineBytes), "\x00")

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid := entry.Name()
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}

		cmdlinePath := filepath.Join("/proc", pid, "cmdline")
		cmdlineBytes, err := os.ReadFile(cmdlinePath)
		if err != nil {
			continue
		}

		cmdlineRaw := strings.TrimRight(string(cmdlineBytes), "\x00")
		if cmdlineRaw == "" || cmdlineRaw == selfCmdlineRaw {
			continue
		}

		cmdline := strings.Join(strings.Split(cmdlineRaw, "\x00"), " ")
		Logger.Debug().Msgf("cmdline: %s", cmdline)

		// Extra exclusion of self process
		// Needed for corner cases
		if strings.Contains(cmdline, "--process "+processName) || strings.Contains(cmdline, "--process="+processName) {
			continue
		}

		if insensitive {
			if sul.ContainsI(cmdlineRaw, processName) {
				return true, nil
			}
		} else {
			if strings.Contains(cmdlineRaw, processName) {
				return true, nil
			}
		}
	}

	return false, nil
}

type ProcessHealthCheck struct {
	ProcessName string
	Insensitive bool
}

func (h *ProcessHealthCheck) Execute() (bool, error) {
	return isProcessRunning(h.ProcessName, h.Insensitive)
}

func (h *ProcessHealthCheck) Dump() string {
	return fmt.Sprintf("process healthcheck for '%s'", h.ProcessName)
}

func (h *ProcessHealthCheck) AreParametersValid() (bool, []string) {
	var errors []string
	if h.ProcessName == "" {
		errors = append(errors, "process name is required")
	}
	return len(errors) == 0, errors
}
