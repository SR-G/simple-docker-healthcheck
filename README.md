

## Cookbooks

### Commands

```bash
simple-docker-healthcheck - single/standalone binary for performing healthchecks in Docker containers without the need for a full Docker image with multiple tools included. It supports various types of healthchecks, including port checks, HTTP status code checks, HTTP response text checks, and HTTP JSON value checks. Replacement of curl, wget, netstat, nc, ..., especially if not available in the container image.

  Usage:
    simple-docker-healthcheck [check-port|check-http-code|check-http-text|check-http-json|check-url]

  Subcommands:
    check-port        Healthcheck that checks if a specific port is open on a host
    check-http-code   Healthcheck that checks the HTTP status code of a specific URL
    check-http-text   Healthcheck that checks if a specific text is present in the HTTP response body of a specific URL
    check-http-json   Healthcheck that checks if a specific JSON value is present in the HTTP response body of a specific URL
    check-url         Healthcheck that checks if a specific URL is reachable (HTTP status code 200-399)

  Flags:
        --version     Displays the program version string.
    -h  --help        Displays help with available flag, subcommand, and positional value parameters.
    -s  --silent      disable all logging output
    -j  --json-logs   enable JSON formatted logs
```

Possible error codes : 

| Error Code | Meaning                                                         |
| ---------- | --------------------------------------------------------------- |
| 0          | Healthcheck has been executed and has passed                    |
| 1          | Healthcheck has been executed but has NOT passed                |
| 2          | Technical error during execution                                |
| 3          | Usage error (missing mandatory parameter for one command, etc.) |


### Check that one port is accessible

```bash
go run . check-port --port 32400
go run . check-port --hostname localhost --port 32400

go run . check-port --hostname mysql --port 3306
```

### Check that a remote URL returns a specific HTTP Status Code

```bash
go run . check-http-code --url "https://gluetun-torrents-exporter.tensin.ovh/metrics" --status-code 200
```

### Check that a remote URL returns an HTTP Status Code included in a specific range

```bash
go run . check-http-code --url "https://gluetun-torrents-exporter.tensin.ovh/metrics" 
--min-status-code 0 --max-status-code 399
```

### Check that a remote URL contains a specific test

```bash
go run . check-http-text --url "https://gluetun-torrents-exporter.tensin.ovh/metrics" --text "go_gc_duration_seconds_sum"
go run . check-http-text --url "https://gluetun-torrents-exporter.tensin.ovh/metrics" --text "go_gc_duration_seconds_sum" --insensitive
```

### Check that a remote URL responding a JSON content has an expected value in a given JSONPath

```bash
go run . check-http-json --url "https://time.tensin.org/" --json-path '$.date' --value '2026-05-16'
```

See JSONPath documentation : 

```bash
0:00 sergio@moon ~/workspaces/projects/simple-docker-healthcheck% curl https://time.tensin.org | jq 
  % Total    % Received % Xferd  Average Speed  Time    Time    Time   Current
                                 Dload  Upload  Total   Spent   Left   Speed
100     73 100     73   0      0   1222      0                              0
{
  "date": "2026-05-16",
  "time": "00:00:50",
  "datetime": "2026-05-16 00:00:50"
}
```

### Docker integration

```dockerfile
COPY --from=ghcr.io/bratteng/healthcheck:latest /sdh /sdh

HEALTHCHECK --interval=5s --timeout=10s --retries=3 CMD [ "/sdh", "--port", "8080" ]
```

### Command mappings between legacy healthchecks and SDH usage

| Instead of ...                                                            | Use ...                                                                                 |
| ------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `curl -f http://localhost:8334/about \|\| exit 1`                         | `/sdh check-url --url http://localhost:8334/about`                                      |
| `curl localhost:3000/api/health 2>/dev/null \| grep -i -q "ok"`           | `/sdh check-http-text --url localhost:3000/api/health --text "ok" -i`                   |
| `bash -c "echo -n '' > /dev/tcp/127.0.0.1/8080"`                          | `/sdh check-port --hostname 127.0.01 --port 9000`                                       |
| `netstat -ltn 2>/dev/null \| grep -c 9000`                                | `/sdh check-port --port 9000`                                                           |
| `nc -z localhost 80`                                                      | `/sdh check-port --hostname localhost --port 80`                                        |
| `curl -f http://localhost/status \| jq '.status' \| grep -i -q "RUNNING"` | `/sdh check-http-json http://localhost/status --json-path '$.status' --value "RUNNING"` |



## Build

### Build the binary

```bash
$ CGO_ENABLED=0 go build -ldflags "-s -w" .
$ upx -9 ./simple-docker-healthcheck
```

### Build the docker image

```bash
docker build . -t simple-docker-healthcheck
(...)
docker run -it --rm --name "test" simple-docker-healthcheck --version
```

### About binary compression

Expected sizes / gains : 

| Binary                                             | Ratio | Size   |
| -------------------------------------------------- | ----- | ------ |
| Built with go compiler                             | 100%  | 13.3 M |
| Built with go compiler + upx -9                    | 56%   | 7.0 M  |
| Built with go compiler with ldflags -s -w          |       | 8.3 M  |
| Built with go compiler with ldflags -s -w + upx -9 |       | 3.3 M  |

At docker image level, then last flavor fives : 

```
IMAGE                              ID             DISK USAGE   CONTENT SIZE   EXTRA
node:24-alpine                     edd927012c1e        160MB             0B        
simple-docker-healthcheck:latest   96ea9531359e       5.46MB             0B        
```


## Links

- https://github.com/bratteng/docker-healthcheck
- https://github.com/bratteng/docker-nginx/blob/15ddec93d6a47ca04f84cdf3bde8b834dee1b806/Dockerfile#L177-L178