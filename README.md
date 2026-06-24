# `simple-docker-healthcheck` (`sdh`)

`simple-docker-healthcheck` is a standalone binary, written in GOLANG, allowing to ease and homogeneize the writing of DOCKER `HEALTHECHECK` commands. Also works with distroless images (i.e., even if curl, netstat, nc, ... are not available in the image, allowing to keep image size under control and maintain a good level of security).

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

Note : if the project is named `simple-docker-healthcheck` (i.e., a self-explanatory name), the binary has been renamed to `sdh` (for the sake of brevity).


### Check that one port is accessible

```bash
sdh check-port --port 32400
sdh check-port --hostname localhost --port 32400

sdh check-port --hostname mysql --port 3306
```

### Check that a remote URL returns a specific HTTP Status Code

```bash
sdh check-http-code --url "https://www.wikipedia.org/" --status-code 200
```

### Check that a remote URL returns an HTTP Status Code included in a specific range

```bash
sdh check-http-code --url "https://www.wikipedia.org/" 
--min-status-code 0 --max-status-code 399
```

### Check that a remote URL contains a specific test

```bash
sdh check-http-text --url "https://www.wikipedia.org/" --text "encyclopedia"
sdh check-http-text --url "https://www.wikipedia.org/" --text "EnCyCloPediA" --insensitive
```

### Check that a remote URL responding a JSON content has an expected value in a given JSONPath

```bash
sdh check-http-json --url "https://time.tensin.org/" --json-path '$.date' --value '2026-05-16'
```

See JSONPath documentation : https://goessner.net/articles/JsonPath/

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
COPY --from=ghcr.io/sr-g/simple-docker-healthcheck:latest /sdh /sdh

HEALTHCHECK --interval=5s --timeout=10s --retries=3 CMD [ "/sdh", "--port", "8080" ]
```

### Command mappings between legacy healthchecks and SDH usage

| Instead of ...                                                            | Use ...                                                                                 |
| ------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `curl -f http://localhost:8334/about \|\| exit 1`                         | `/sdh check-url --url http://localhost:8334/about`                                      |
| `curl localhost:3000/api/health 2>/dev/null \| grep -i -q "ok"`           | `/sdh check-http-text --url localhost:3000/api/health --text "ok" -i`                   |
| `bash -c "echo -n '' > /dev/tcp/127.0.0.1/8080"`                          | `/sdh check-port --hostname 127.0.0.1 --port 9000`                                      |
| `netstat -ltn 2>/dev/null \| grep -c 9000`                                | `/sdh check-port --port 9000`                                                           |
| `nc -z localhost 80`                                                      | `/sdh check-port --hostname localhost --port 80`                                        |
| `curl -f http://localhost/status \| jq '.status' \| grep -i -q "RUNNING"` | `/sdh check-http-json http://localhost/status --json-path '$.status' --value "RUNNING"` |



## Build

### Build the binary

```bash
$ CGO_ENABLED=0 go build -ldflags "-s -w" .
$ upx -9 ./simple-docker-healthcheck
```

Or use the makefile : 

```bash
make build
make docker-build
make docker-run

### Build the docker image

```bash
docker build . -t simple-docker-healthcheck
(...)
docker run -it --rm --name "test" simple-docker-healthcheck --version
```

### About binary compression

Expected sizes / gains when using various flags and/or UPX : 

| Binary                                             | Ratio | Size   |
| -------------------------------------------------- | ----- | ------ |
| Built with go compiler                             | 100%  | 13.3 M |
| Built with go compiler + upx -9                    | 56%   | 7.0 M  |
| Built with go compiler with ldflags -s -w          |       | 8.3 M  |
| Built with go compiler with ldflags -s -w + upx -9 |       | 3.3 M  |

At docker image level, then last flavor gives : 

```
IMAGE                              ID             DISK USAGE   CONTENT SIZE   EXTRA
simple-docker-healthcheck:latest   96ea9531359e       5.46MB             0B        
```

Note : this is just a preliminary investigation, for now the binaries are not released with UPX being used.


## Links

- `docker-healthcheck` : a similar preliminary implementation (https://github.com/bratteng/docker-healthcheck) (archived)
- usage of that `docker-healthcheck` binary in the NGINX docker image, directly from the docker container, as an example of usage (https://github.com/bratteng/docker-nginx/blob/15ddec93d6a47ca04f84cdf3bde8b834dee1b806/Dockerfile#L177-L178)
- https://itnext.io/healthchecks-with-distroless-containers-262a52abc31e