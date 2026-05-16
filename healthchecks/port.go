package healthchecks

import (
	"fmt"
	"net"
	"simple-docker-healthcheck/constants"
)

type PortHealthCheck struct {
	Hostname string
	Port     int
}

func (p *PortHealthCheck) Execute() (bool, error) {
	addr := net.JoinHostPort(p.Hostname, fmt.Sprintf("%d", p.Port))
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false, err
	}
	conn.Close()
	return true, nil
}

func (p *PortHealthCheck) Dump() string {
	return fmt.Sprintf("port healthcheck on hostname '%s', port '%d'}", p.Hostname, p.Port)
}

func (p *PortHealthCheck) AreParametersValid() (bool, []string) {
	var errors []string
	if p.Hostname == "" {
		errors = append(errors, "hostname is required")
	}
	if p.Port == constants.HTTP_STATUS_CODE_UNSET {
		errors = append(errors, "port is required")
	}
	return len(errors) == 0, errors
}
