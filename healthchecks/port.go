package healthchecks

import (
	"fmt"
	"net"
)

type PortHealthCheck struct {
	Hostname string
	Port     int
}

func (p *PortHealthCheck) Execute() (bool, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", p.Hostname, p.Port))
	if err != nil {
		return false, err
	}
	conn.Close()
	return true, nil
}
