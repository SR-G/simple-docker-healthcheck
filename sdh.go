package main

import (
	"fmt"
	"os"
	"simple-docker-healthcheck/constants"
	"simple-docker-healthcheck/healthchecks"
	"time"

	"github.com/integrii/flaggy"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger
var osExit = os.Exit

const (
	PROGRAM_NAME          = "sdh"
	PROGRAM_VERSION       = "1.0.0"    // overwritten at release time from makefile
	PROGRAM_VERSION_LABEL = "SNAPSHOT" // overwritten at release time from makefile
)

func newLogger(jsonLogs, debug, silent bool) zerolog.Logger {
	if silent {
		return zerolog.Nop()
	}

	var l zerolog.Logger
	if jsonLogs {
		l = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true}
		l = zerolog.New(writer).With().Timestamp().Logger()
	}
	if debug {
		l = l.Level(zerolog.DebugLevel)
	} else {
		l = l.Level(zerolog.InfoLevel)
	}
	return l
}

func initParametersAndParseFlags() (healthchecks.HealthCheck, error) {
	var logsSilent bool = false
	var logsJson bool = false
	var logsDebug bool = false
	flaggy.Bool(&logsSilent, "s", "silent", "disable all logging output")
	flaggy.Bool(&logsJson, "j", "json-logs", "enable JSON formatted logs")
	flaggy.Bool(&logsDebug, "d", "debug", "enable debug logging")

	var commandPortHealthCheckHostname string = constants.LOCALHOST
	var commandPortHealthCheckPort int = constants.HTTP_STATUS_CODE_UNSET
	commandPortHealthCheck := flaggy.NewSubcommand("check-port")
	commandPortHealthCheck.Description = "Healthcheck that checks if a specific port is open on a host"
	commandPortHealthCheck.String(&commandPortHealthCheckHostname, "", "hostname", "hostname to check")
	commandPortHealthCheck.Int(&commandPortHealthCheckPort, "", "port", "port to check")

	var commandHTTPCodeHealthCheckURL string = ""
	var commandHTTPCodeHealthCheckExpectedStatusCode int = constants.HTTP_STATUS_CODE_OK_MIN
	var commandHTTPCodeHealthCheckExpectedStatusCodeMin int = constants.HTTP_STATUS_CODE_UNSET
	var commandHTTPCodeHealthCheckExpectedStatusCodeMax int = constants.HTTP_STATUS_CODE_UNSET
	commandHTTPCodeHealthCheck := flaggy.NewSubcommand("check-http-code")
	commandHTTPCodeHealthCheck.Description = "Healthcheck that checks the HTTP status code of a specific URL"
	commandHTTPCodeHealthCheck.String(&commandHTTPCodeHealthCheckURL, "", "url", "URL to check")
	commandHTTPCodeHealthCheck.Int(&commandHTTPCodeHealthCheckExpectedStatusCode, "", "status-code", "expected HTTP status code")
	commandHTTPCodeHealthCheck.Int(&commandHTTPCodeHealthCheckExpectedStatusCodeMin, "", "min-status-code", "expected minimum HTTP status code (ranged healthcheck)")
	commandHTTPCodeHealthCheck.Int(&commandHTTPCodeHealthCheckExpectedStatusCodeMax, "", "max-status-code", "expected maximum HTTP status code (ranged healthcheck)")

	var commandHTTPTextHealthCheckURL string = ""
	var commandHTTPTextHealthCheckExpectedText string = ""
	var commandHTTPTextHealthCheckInsensitive bool = false
	commandHTTPTextHealthCheck := flaggy.NewSubcommand("check-http-text")
	commandHTTPTextHealthCheck.Description = "Healthcheck that checks if a specific text is present in the HTTP response body of a specific URL"
	commandHTTPTextHealthCheck.String(&commandHTTPTextHealthCheckURL, "", "url", "URL to check")
	commandHTTPTextHealthCheck.String(&commandHTTPTextHealthCheckExpectedText, "", "text", "expected text in the HTTP response body")
	commandHTTPTextHealthCheck.Bool(&commandHTTPTextHealthCheckInsensitive, "i", "insensitive", "perform case-insensitive text matching")

	var commandHTTPJSONHealthCheckURL string = ""
	var commandHTTPJSONHealthCheckJSONPath string = ""
	var commandHTTPJSONHealthCheckExpectedValue string = ""
	var commandHTTPJSONHealthCheckInsensitive bool = false
	commandHTTPJSONHealthCheck := flaggy.NewSubcommand("check-http-json")
	commandHTTPJSONHealthCheck.Description = "Healthcheck that checks if a specific JSON value is present in the HTTP response body of a specific URL"
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckURL, "", "url", "URL to check")
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckJSONPath, "", "json-path", "JSONPath of the value to check in the HTTP response body")
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckExpectedValue, "", "value", "expected value at the specified JSONPath in the HTTP response body")
	commandHTTPJSONHealthCheck.Bool(&commandHTTPJSONHealthCheckInsensitive, "i", "insensitive", "perform case-insensitive JSON value matching")

	var commandURLCheckURL string = ""
	commandURLCheck := flaggy.NewSubcommand("check-url")
	commandURLCheck.Description = "Healthcheck that checks if a specific URL is reachable (HTTP status code 200-399)"
	commandURLCheck.String(&commandURLCheckURL, "", "url", "URL to check")

	var commandProcessHealthCheckProcessName string = ""
	commandProcessHealthCheck := flaggy.NewSubcommand("check-process")
	commandProcessHealthCheck.Description = "Healthcheck that checks if a specific process is running (linux only)"
	commandProcessHealthCheck.String(&commandProcessHealthCheckProcessName, "", "process", "name of the process to check")

	flaggy.AttachSubcommand(commandPortHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPCodeHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPTextHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPJSONHealthCheck, 1)
	flaggy.AttachSubcommand(commandURLCheck, 1)
	flaggy.AttachSubcommand(commandProcessHealthCheck, 1)

	flaggy.SetDescription("single/standalone binary for performing healthchecks in Docker containers without the need for a full Docker image with multiple tools included. It supports various types of healthchecks, including port checks, HTTP status code checks, HTTP response text checks, and HTTP JSON value checks. Replacement of curl, wget, netstat, nc, ..., especially if not available in the container image.")
	flaggy.SetVersion(Version.String())
	flaggy.DisableCompletion()
	flaggy.Parse()

	logger = newLogger(logsJson, logsDebug, logsSilent)
	healthchecks.Logger = logger

	switch {
	case commandPortHealthCheck.Used:
		return &healthchecks.PortHealthCheck{
			Hostname: commandPortHealthCheckHostname,
			Port:     commandPortHealthCheckPort,
		}, nil
	case commandHTTPCodeHealthCheck.Used:
		if commandHTTPCodeHealthCheckExpectedStatusCodeMin != constants.HTTP_STATUS_CODE_UNSET && commandHTTPCodeHealthCheckExpectedStatusCodeMax != constants.HTTP_STATUS_CODE_UNSET {
			return &healthchecks.HTTPCodeHealthCheckRangeStatusCode{
				URL:                   commandHTTPCodeHealthCheckURL,
				MinExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCodeMin,
				MaxExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCodeMax,
			}, nil
		} else {
			return &healthchecks.HTTPCodeHealthCheckExactStatusCode{
				URL:                commandHTTPCodeHealthCheckURL,
				ExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCode,
			}, nil
		}
	case commandURLCheck.Used:
		return &healthchecks.HTTPCodeHealthCheckRangeStatusCode{
			URL:                   commandURLCheckURL,
			MinExpectedStatusCode: constants.HTTP_STATUS_CODE_OK_MIN,
			MaxExpectedStatusCode: constants.HTTP_STATUS_CODE_OK_MAX,
		}, nil
	case commandHTTPTextHealthCheck.Used:
		return &healthchecks.HTTPHealthCheckText{
			URL:          commandHTTPTextHealthCheckURL,
			ExpectedText: commandHTTPTextHealthCheckExpectedText,
			Insensitive:  commandHTTPTextHealthCheckInsensitive,
		}, nil
	case commandHTTPJSONHealthCheck.Used:
		return &healthchecks.HTTPHealthCheckJSONPath{
			URL:           commandHTTPJSONHealthCheckURL,
			JSONPath:      commandHTTPJSONHealthCheckJSONPath,
			ExpectedValue: commandHTTPJSONHealthCheckExpectedValue,
			Insensitive:   commandHTTPJSONHealthCheckInsensitive,
		}, nil
	case commandProcessHealthCheck.Used:
		return &healthchecks.ProcessHealthCheck{
			ProcessName: commandProcessHealthCheckProcessName,
		}, nil
	default:
		return nil, fmt.Errorf("no subcommand provided")
	}
}

func executeHealthCheck(health healthchecks.HealthCheck) {
	if valid, validationErrors := health.AreParametersValid(); !valid {
		logger.Error().Str("healthcheck", health.Dump()).Interface("validation_errors", validationErrors).Msg("invalid parameters")
		osExit(constants.OS_EXIT_CODE_USAGE)
	}

	logger.Info().Str("healthcheck", health.Dump()).Msg("executing")
	ok, err := health.Execute()
	if err != nil {
		logger.Error().Err(err).Msg("healthcheck execution failed")
		osExit(constants.OS_EXIT_CODE_ERROR)
	}

	if ok {
		logger.Info().Msg("healthcheck passed")
		osExit(constants.OS_EXIT_CODE_SUCCESS)
	} else {
		logger.Warn().Msg("healthcheck failed")
		osExit(constants.OS_EXIT_CODE_FAILURE)
	}
}

func main() {
	healthCheck, err := initParametersAndParseFlags()
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize healthcheck parameters")
		osExit(constants.OS_EXIT_CODE_USAGE)
	}
	executeHealthCheck(healthCheck)
}
