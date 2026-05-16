package healthchecks

type HealthCheck interface {
	Execute() (bool, error)
	Dump() string
	AreParametersValid() (bool, []string)
}
