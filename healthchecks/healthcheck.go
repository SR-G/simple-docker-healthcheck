package healthchecks

type HealthCheck interface {
	Execute() (bool, error)
}
