package main

import (
	"os"
	"simple-docker-healthcheck/constants"
	"simple-docker-healthcheck/healthchecks"
	"time"

	"github.com/integrii/flaggy"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger
var osExit = os.Exit

func main() {

	var silent bool = false
	var jsonLogs bool = false
	flaggy.Bool(&silent, "s", "silent", "disable all logging output")
	flaggy.Bool(&jsonLogs, "j", "json-logs", "enable JSON formatted logs")

	var commandPortHealthCheckHostname string = "localhost"
	var commandPortHealthCheckPort int
	commandPortHealthCheck := flaggy.NewSubcommand("check-port")
	commandPortHealthCheck.Description = "Healthcheck that checks if a specific port is open on a host"
	commandPortHealthCheck.String(&commandPortHealthCheckHostname, "", "hostname", "hostname to check")
	commandPortHealthCheck.Int(&commandPortHealthCheckPort, "", "port", "port to check")

	var commandHTTPCodeHealthCheckURL string
	var commandHTTPCodeHealthCheckExpectedStatusCode int = 200
	var commandHTTPCodeHealthCheckExpectedStatusCodeMin int
	var commandHTTPCodeHealthCheckExpectedStatusCodeMax int
	commandHTTPCodeHealthCheck := flaggy.NewSubcommand("check-http-code")
	commandHTTPCodeHealthCheck.Description = "Healthcheck that checks the HTTP status code of a specific URL"
	commandHTTPCodeHealthCheck.String(&commandHTTPCodeHealthCheckURL, "", "url", "URL to check")
	commandHTTPCodeHealthCheck.Int(&commandHTTPCodeHealthCheckExpectedStatusCode, "", "status-code", "expected HTTP status code")
	commandHTTPCodeHealthCheck.Int(&commandHTTPCodeHealthCheckExpectedStatusCodeMin, "", "min-status-code", "expected minimum HTTP status code (ranged healthcheck)")
	commandHTTPCodeHealthCheck.Int(&commandHTTPCodeHealthCheckExpectedStatusCodeMax, "", "max-status-code", "expected maximum HTTP status code (ranged healthcheck)")

	var commandHTTPTextHealthCheckURL string
	var commandHTTPTextHealthCheckExpectedText string
	var commandHTTPTextHealthCheckInsensitive bool
	commandHTTPTextHealthCheck := flaggy.NewSubcommand("check-http-text")
	commandHTTPTextHealthCheck.Description = "Healthcheck that checks if a specific text is present in the HTTP response body of a specific URL"
	commandHTTPTextHealthCheck.String(&commandHTTPTextHealthCheckURL, "", "url", "URL to check")
	commandHTTPTextHealthCheck.String(&commandHTTPTextHealthCheckExpectedText, "", "text", "expected text in the HTTP response body")
	commandHTTPTextHealthCheck.Bool(&commandHTTPTextHealthCheckInsensitive, "i", "insensitive", "perform case-insensitive text matching")

	var commandHTTPJSONHealthCheckURL string
	var commandHTTPJSONHealthCheckJSONPath string
	var commandHTTPJSONHealthCheckExpectedValue string
	var commandHTTPJSONHealthCheckInsensitive bool
	commandHTTPJSONHealthCheck := flaggy.NewSubcommand("check-http-json")
	commandHTTPJSONHealthCheck.Description = "Healthcheck that checks if a specific JSON value is present in the HTTP response body of a specific URL"
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckURL, "", "url", "URL to check")
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckJSONPath, "", "json-path", "JSONPath of the value to check in the HTTP response body")
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckExpectedValue, "", "value", "expected value at the specified JSONPath in the HTTP response body")
	commandHTTPJSONHealthCheck.Bool(&commandHTTPJSONHealthCheckInsensitive, "i", "insensitive", "perform case-insensitive JSON value matching")
	commandURLCheck := flaggy.NewSubcommand("check-url")
	commandURLCheck.Description = "Healthcheck that checks if a specific URL is reachable (HTTP status code 200-399)"
	var commandURLCheckURL string
	commandURLCheck.String(&commandURLCheckURL, "", "url", "URL to check")

	flaggy.AttachSubcommand(commandPortHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPCodeHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPTextHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPJSONHealthCheck, 1)
	flaggy.AttachSubcommand(commandURLCheck, 1)

	flaggy.SetDescription("single/standalone binary for performing healthchecks in Docker containers without the need for a full Docker image with multiple tools included. It supports various types of healthchecks, including port checks, HTTP status code checks, HTTP response text checks, and HTTP JSON value checks. Replacement of curl, wget, netstat, nc, ..., especially if not available in the container image.")
	flaggy.SetVersion("1.0.0")
	flaggy.DisableCompletion()
	flaggy.Parse()

	if jsonLogs {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}

	if silent {
		logger = logger.Level(zerolog.Disabled)
	}

	if commandPortHealthCheck != nil && commandPortHealthCheck.Used {

		logger.Info().Msgf("executing port healthcheck on %s:%d", commandPortHealthCheckHostname, commandPortHealthCheckPort)
		healthCheck := &healthchecks.PortHealthCheck{
			Hostname: commandPortHealthCheckHostname,
			Port:     commandPortHealthCheckPort,
		}

		executeHealthCheck(healthCheck)
	} else if commandHTTPCodeHealthCheck != nil && commandHTTPCodeHealthCheck.Used {

		if commandHTTPCodeHealthCheckExpectedStatusCodeMin >= 0 && commandHTTPCodeHealthCheckExpectedStatusCodeMax > 0 {
			logger.Info().Msgf("executing HTTP code healthcheck on %s with expected status code range %d-%d", commandHTTPCodeHealthCheckURL, commandHTTPCodeHealthCheckExpectedStatusCodeMin, commandHTTPCodeHealthCheckExpectedStatusCodeMax)
			healthCheck := &healthchecks.HTTPCodeHealthCheckRangeStatusCode{
				URL:                   commandHTTPCodeHealthCheckURL,
				MinExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCodeMin,
				MaxExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCodeMax,
			}

			executeHealthCheck(healthCheck)
		} else {
			logger.Info().Msgf("executing HTTP code healthcheck on %s with expected status code %d", commandHTTPCodeHealthCheckURL, commandHTTPCodeHealthCheckExpectedStatusCode)
			healthCheck := &healthchecks.HTTPCodeHealthCheckExactStatusCode{
				URL:                commandHTTPCodeHealthCheckURL,
				ExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCode,
			}

			executeHealthCheck(healthCheck)
		}
	} else if commandURLCheck != nil && commandURLCheck.Used {
		logger.Info().Msgf("executing URL healthcheck on %s", commandURLCheckURL)
		healthCheck := &healthchecks.HTTPCodeHealthCheckRangeStatusCode{
			URL:                   commandURLCheckURL,
			MinExpectedStatusCode: 200,
			MaxExpectedStatusCode: 399,
		}
		executeHealthCheck(healthCheck)
	} else if commandHTTPTextHealthCheck != nil && commandHTTPTextHealthCheck.Used {

		logger.Info().Msgf("executing HTTP text healthcheck on %s with expected text '%s'", commandHTTPTextHealthCheckURL, commandHTTPTextHealthCheckExpectedText)
		healthCheck := &healthchecks.HTTPHealthCheckText{
			URL:          commandHTTPTextHealthCheckURL,
			ExpectedText: commandHTTPTextHealthCheckExpectedText,
			Insensitive:  commandHTTPTextHealthCheckInsensitive,
		}

		executeHealthCheck(healthCheck)
	} else if commandHTTPJSONHealthCheck != nil && commandHTTPJSONHealthCheck.Used {

		logger.Info().Msgf("executing HTTP JSON healthcheck on %s with expected value '%s' at JSONPath '%s'", commandHTTPJSONHealthCheckURL, commandHTTPJSONHealthCheckExpectedValue, commandHTTPJSONHealthCheckJSONPath)
		healthCheck := &healthchecks.HTTPHealthCheckJSONPath{
			URL:           commandHTTPJSONHealthCheckURL,
			JSONPath:      commandHTTPJSONHealthCheckJSONPath,
			ExpectedValue: commandHTTPJSONHealthCheckExpectedValue,
			Insensitive:   commandHTTPJSONHealthCheckInsensitive,
		}

		executeHealthCheck(healthCheck)
	} else {
		logger.Error().Msg("no subcommand provided")
		osExit(constants.OS_EXIT_CODE_USAGE)
	}
}

func executeHealthCheck(healthCheck healthchecks.HealthCheck) {
	ok, err := healthCheck.Execute()
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
