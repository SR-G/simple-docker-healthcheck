package main

import (
	"os"
	"simple-docker-healthcheck/healthchecks"

	"github.com/integrii/flaggy"
	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func main() {

	var silent bool = false
	flaggy.Bool(&silent, "s", "silent", "disable all logging output")

	var commandPortHealthCheckHostname string = "localhost"
	var commandPortHealthCheckPort int
	commandPortHealthCheck := flaggy.NewSubcommand("check-port")
	commandPortHealthCheck.Description = "Healthcheck that checks if a specific port is open on a host"
	commandPortHealthCheck.String(&commandPortHealthCheckHostname, "", "hostname", "hostname to check")
	commandPortHealthCheck.Int(&commandPortHealthCheckPort, "", "port", "port to check")

	var commandHTTPCodeHealthCheckURL string
	var commandHTTPCodeHealthCheckExpectedStatusCode int
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
	commandHTTPTextHealthCheck := flaggy.NewSubcommand("check-http-text")
	commandHTTPTextHealthCheck.Description = "Healthcheck that checks if a specific text is present in the HTTP response body of a specific URL"
	commandHTTPTextHealthCheck.String(&commandHTTPTextHealthCheckURL, "", "url", "URL to check")
	commandHTTPTextHealthCheck.String(&commandHTTPTextHealthCheckExpectedText, "", "text", "expected text in the HTTP response body")

	var commandHTTPJSONHealthCheckURL string
	var commandHTTPJSONHealthCheckJSONPath string
	var commandHTTPJSONHealthCheckExpectedValue string
	commandHTTPJSONHealthCheck := flaggy.NewSubcommand("check-http-json")
	commandHTTPJSONHealthCheck.Description = "Healthcheck that checks if a specific JSON value is present in the HTTP response body of a specific URL"
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckURL, "", "url", "URL to check")
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckJSONPath, "", "json-path", "JSONPath of the value to check in the HTTP response body")
	commandHTTPJSONHealthCheck.String(&commandHTTPJSONHealthCheckExpectedValue, "", "value", "expected value at the specified JSONPath in the HTTP response body")

	flaggy.AttachSubcommand(commandPortHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPCodeHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPTextHealthCheck, 1)
	flaggy.AttachSubcommand(commandHTTPJSONHealthCheck, 1)

	flaggy.SetDescription("single/standalone binary for performing healthchecks in Docker containers without the need for a full Docker image with multiple tools included. It supports various types of healthchecks, including port checks, HTTP status code checks, HTTP response text checks, and HTTP JSON value checks. Replacement of curl, wget, netstat, nc, ..., especially if not available in the container image.")
	flaggy.SetVersion("1.0.0")
	flaggy.Parse()

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

		if commandHTTPCodeHealthCheckExpectedStatusCode != 0 {
			logger.Info().Msgf("executing HTTP code healthcheck on %s with expected status code %d", commandHTTPCodeHealthCheckURL, commandHTTPCodeHealthCheckExpectedStatusCode)
			healthCheck := &healthchecks.HTTPCodeHealthCheckExactStatusCode{
				URL:                commandHTTPCodeHealthCheckURL,
				ExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCode,
			}

			executeHealthCheck(healthCheck)
		} else if commandHTTPCodeHealthCheckExpectedStatusCodeMin >= 0 && commandHTTPCodeHealthCheckExpectedStatusCodeMax >= 0 {
			logger.Info().Msgf("executing HTTP code healthcheck on %s with expected status code range %d-%d", commandHTTPCodeHealthCheckURL, commandHTTPCodeHealthCheckExpectedStatusCodeMin, commandHTTPCodeHealthCheckExpectedStatusCodeMax)
			healthCheck := &healthchecks.HTTPCodeHealthCheckRangeStatusCode{
				URL:                   commandHTTPCodeHealthCheckURL,
				MinExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCodeMin,
				MaxExpectedStatusCode: commandHTTPCodeHealthCheckExpectedStatusCodeMax,
			}

			executeHealthCheck(healthCheck)
		} else {
			logger.Error().Msg("no expected status code provided for HTTP code healthcheck")
			os.Exit(1)
		}
	} else if commandHTTPTextHealthCheck != nil && commandHTTPTextHealthCheck.Used {

		logger.Info().Msgf("executing HTTP text healthcheck on %s with expected text '%s'", commandHTTPTextHealthCheckURL, commandHTTPTextHealthCheckExpectedText)
		healthCheck := &healthchecks.HTTPHealthCheckText{
			URL:          commandHTTPTextHealthCheckURL,
			ExpectedText: commandHTTPTextHealthCheckExpectedText,
		}

		executeHealthCheck(healthCheck)
	} else if commandHTTPJSONHealthCheck != nil && commandHTTPJSONHealthCheck.Used {

		logger.Info().Msgf("executing HTTP JSON healthcheck on %s with expected value '%s' at JSONPath '%s'", commandHTTPJSONHealthCheckURL, commandHTTPJSONHealthCheckExpectedValue, commandHTTPJSONHealthCheckJSONPath)
		healthCheck := &healthchecks.HTTPHealthCheckJSONPath{
			URL:           commandHTTPJSONHealthCheckURL,
			JSONPath:      commandHTTPJSONHealthCheckJSONPath,
			ExpectedValue: commandHTTPJSONHealthCheckExpectedValue,
		}

		executeHealthCheck(healthCheck)
	} else {
		logger.Error().Msg("no subcommand provided")
		os.Exit(1)
	}
}

func executeHealthCheck(healthCheck healthchecks.HealthCheck) {
	ok, err := healthCheck.Execute()
	if err != nil {
		logger.Error().Err(err).Msg("healthcheck execution failed")
		os.Exit(1)
	}

	if ok {
		logger.Info().Msg("healthcheck passed")
		os.Exit(0)
	} else {
		logger.Warn().Msg("healthcheck failed")
		os.Exit(1)
	}
}
