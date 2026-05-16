package constants

const (
	DEFAULT_USER_AGENT = "Simple-Docker-Healthcheck/1.0.0"

	OS_EXIT_CODE_SUCCESS = 0 // healthcheck passed
	OS_EXIT_CODE_FAILURE = 1 // healthcheck failed
	OS_EXIT_CODE_ERROR   = 2 // error during healthcheck execution
	OS_EXIT_CODE_USAGE   = 3 // incorrect usage (e.g., missing subcommand or parameters)

	HTTP_STATUS_CODE_OK_MIN = 200
	HTTP_STATUS_CODE_OK_MAX = 399
	HTTP_STATUS_CODE_UNSET  = -1

	LOCALHOST = "localhost"
)
