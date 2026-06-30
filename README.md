# `simple-docker-healthcheck` (`sdh`)

- [`simple-docker-healthcheck` (`sdh`)](#simple-docker-healthcheck-sdh)
  - [Cookbooks](#cookbooks)
    - [Commands (flags)](#commands-flags)
    - [Recipe - Check that one port is accessible](#recipe---check-that-one-port-is-accessible)
    - [Recipe - Check that a remote URL returns a specific HTTP Status Code](#recipe---check-that-a-remote-url-returns-a-specific-http-status-code)
    - [Recipe - Check that a remote URL returns an HTTP Status Code included in a specific range](#recipe---check-that-a-remote-url-returns-an-http-status-code-included-in-a-specific-range)
    - [Recipe - Check that a remote URL contains a specific text](#recipe---check-that-a-remote-url-contains-a-specific-text)
    - [Recipe - Check that a remote URL responding a JSON content has an expected value in a given JSONPath](#recipe---check-that-a-remote-url-responding-a-json-content-has-an-expected-value-in-a-given-jsonpath)
    - [Recipe - Check that a file is available on filesystem](#recipe---check-that-a-file-is-available-on-filesystem)
    - [Recipe - Check that a file is available on filesystem and has a specific content](#recipe---check-that-a-file-is-available-on-filesystem-and-has-a-specific-content)
    - [Recipe - Check that a file is available on filesystem and is matching a REGEXP](#recipe---check-that-a-file-is-available-on-filesystem-and-is-matching-a-regexp)
    - [Recipe - Check that a process is available in memory](#recipe---check-that-a-process-is-available-in-memory)
    - [Docker integration](#docker-integration)
    - [Command mappings between legacy healthchecks and `sdh` usage](#command-mappings-between-legacy-healthchecks-and-sdh-usage)
    - [Example of docker integrations](#example-of-docker-integrations)
  - [DEV Activities](#dev-activities)
    - [Build](#build)
      - [Build the binary](#build-the-binary)
      - [Build the docker image](#build-the-docker-image)
      - [Release](#release)
    - [Overwrite a previous tag](#overwrite-a-previous-tag)
    - [About binary compression](#about-binary-compression)
  - [Links](#links)

`simple-docker-healthcheck` is a standalone binary, written in GOLANG, allowing to ease and homogeneize the writing of DOCKER `HEALTHECHECK` commands. 

This allows to use the exact same binary for multiple kind of healthchecks : 

- **Checking if a port is opened at TCP level** (replaces `netstat`, `nc`)
- **Checking if a URL is reachable** (replaces `wget`, `curl`)
- **Checking the HTTP code of a URL** (either a specific value, either a range) (replaces `wget`, `curl`)
- **Checking is a specific TEXT string is found in the response of a URL call** (replaces `wget`, `curl` + `grep`)
- **Checking a JSON-Path expression in the response returned by a URL call** (replaces `wget`, `curl` + `jq`)
- **Checking for a specific process in-memory** (replaces `ps` - a few DOCKER parent images may not be providing a shell and a ps command)

Also works very well with distroless images (i.e., even if curl, netstat, nc, ... are not available in the image), allowing : 

- to keep image size under control (no need to add extra dependencies)
- to maintain a good level of security (only one extra single binary to be added for the healthcheck)

## Cookbooks

You'll find below some examples, and at the end a "transformation table" between "legacy" commands and `sdh` counterparts.

### Commands (flags)

```bash
sdh - single/standalone binary for performing healthchecks in Docker containers without the need for a full Docker image with multiple tools included. It supports various types of healthchecks, including port checks, HTTP status code checks, HTTP response text checks, and HTTP JSON value checks. Replacement of curl, wget, netstat, nc, ..., especially if not available in the container image.

  Usage:
    sdh [check-port|check-http-code|check-http-text|check-http-json|check-url|check-process|check-file|check-file-content|check-file-regexp]

  Subcommands:
    check-port           Healthcheck that checks if a specific port is open on a host
    check-http-code      Healthcheck that checks the HTTP status code of a specific URL
    check-http-text      Healthcheck that checks if a specific text is present in the HTTP response body of a specific URL
    check-http-json      Healthcheck that checks if a specific JSON value is present in the HTTP response body of a specific URL
    check-url            Healthcheck that checks if a specific URL is reachable (HTTP status code 200-399)
    check-process        Healthcheck that checks if a specific process is running (linux only)
    check-file           Healthcheck that checks if a specific file is available on filesystem (like a .PID file)
    check-file-content   Healthcheck that checks if a specific file has a specific content
    check-file-regexp    Healthcheck that checks if a specific file has a specific content, thanks to a regular expression

  Flags:
        --version     Displays the program version string.
    -h  --help        Displays help with available flag, subcommand, and positional value parameters.
    -s  --silent      disable all logging output
    -j  --json-logs   enable JSON formatted logs
    -d  --debug       enable debug logging
```

Possible error codes : 

| Error Code | Meaning                                                         |
| ---------- | --------------------------------------------------------------- |
| 0          | Healthcheck has been executed and has passed                    |
| 1          | Healthcheck has been executed but has NOT passed                |
| 2          | Technical error during execution                                |
| 3          | Usage error (missing mandatory parameter for one command, etc.) |

Notes : 
- if the project is named `simple-docker-healthcheck` (i.e., a self-explanatory name), the binary has been renamed to `sdh` (for the sake of brevity).
- each command may have additional parameters, just use `sdh <command> --help`


### Recipe - Check that one port is accessible

```bash
sdh check-port --port 32400
sdh check-port --hostname localhost --port 32400

sdh check-port --hostname mysql --port 3306
```

### Recipe - Check that a remote URL returns a specific HTTP Status Code

```bash
sdh check-http-code --url "https://www.wikipedia.org/" --status-code 200
```

Example of extra parameters : 

```bash
check-http-code - Healthcheck that checks the HTTP status code of a specific URL

  Flags:
      --url               URL to check
      --status-code       expected HTTP status code (default: 200)
      --min-status-code   expected minimum HTTP status code (ranged healthcheck) (default: -1)
      --max-status-code   expected maximum HTTP status code (ranged healthcheck) (default: -1)
```

### Recipe - Check that a remote URL returns an HTTP Status Code included in a specific range

```bash
sdh check-http-code --url "https://www.wikipedia.org/" 
--min-status-code 0 --max-status-code 399
```

### Recipe - Check that a remote URL contains a specific text

```bash
sdh check-http-text --url "https://www.wikipedia.org/" --text "encyclopedia"
sdh check-http-text --url "https://www.wikipedia.org/" --text "EnCyCloPediA" --insensitive
```

### Recipe - Check that a remote URL responding a JSON content has an expected value in a given JSONPath

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

### Recipe - Check that a file is available on filesystem

```bash
sdh check-file --filename "/var/run/myapp.pid"
```

### Recipe - Check that a file is available on filesystem and has a specific content

```bash
sdh check-file-content --filename "/var/run/myapp.pid" --content "RUNNING" --insensitive
```

### Recipe - Check that a file is available on filesystem and is matching a REGEXP

```bash
# Check that the PID file of myapp is containing a valid digit
sdh check-file-regexp --filename "/var/run/myapp.pid" --regexp '[0-9]+'
```

### Recipe - Check that a process is available in memory

```bash
# Check that the process is available in memory
sdh check-process --process "sdh'
```

Notes : 
- most of the time, this does not make sense (as the container should exit if the main process is exited), but this may be helpful for some corner cases (s6 as the main process and multiple other sub-processes, etc.)
- this is NOT relying on the `ps` command (only on `/proc`)
- as a consequence, at this time, this is only working on linux and is not available with the windows binary


### Docker integration

You can directly embed the binary from the corresponding `sdh` docker image published in `Github Repositories` : https://github.com/SR-G/simple-docker-healthcheck/pkgs/container/simple-docker-healthcheck

```dockerfile
COPY --from=ghcr.io/sr-g/simple-docker-healthcheck:latest /sdh /sdh

HEALTHCHECK --interval=5s --timeout=10s --retries=3 CMD [ "/sdh", "check-port", "--port", "8080" ]
```

### Command mappings between legacy healthchecks and `sdh` usage

| Instead of ...                                                            | Use ...                                                                                 |
| ------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `curl -f http://localhost:8334/about \|\| exit 1`                         | `/sdh check-url --url http://localhost:8334/about`                                      |
| `curl localhost:3000/api/health 2>/dev/null \| grep -i -q "ok"`           | `/sdh check-http-text --url localhost:3000/api/health --text "ok" -i`                   |
| `bash -c "echo -n '' > /dev/tcp/127.0.0.1/8080"`                          | `/sdh check-port --hostname 127.0.0.1 --port 9000`                                      |
| `netstat -ltn 2>/dev/null \| grep -c 9000`                                | `/sdh check-port --port 9000`                                                           |
| `nc -z localhost 80`                                                      | `/sdh check-port --hostname localhost --port 80`                                        |
| `[ ! -f /var/run/myapp.pid ] && exit 1 `                                  | `/sdh check-file --filename /var/run/myapp.pid`  |
| `grep -c -i "RUNNING" /var/run/myapp.pid `                                | `/sdh check-file-content --filename /var/run/myapp.pid --content "RUNNING" --insensitive` |
| `grep -c -e "[0-9]+" /var/run/myapp.pid`                                  | `/sdh check-file-regexp --filename /var/run/myapp.pid --regexp "[0-9]+"` |
| `curl -f http://localhost/status \| jq '.status' \| grep -i -q "RUNNING"` | `/sdh check-http-json http://localhost/status --json-path '$.status' --value "RUNNING"` |
| `grep -ql '[R]oonServer.dll' /proc/[0-9]*/cmdline 2>/dev/null || exit 1`  | `/sdh check-process --process "RoonServer.dll" --insensitive |


### Example of docker integrations

In a Dockerfile (see Docker Reference : https://docs.docker.com/reference/dockerfile/#healthcheck) : 

```
FROM        <my_parent_image>

HEALTHCHECK --interval=60s --start-period=15s CMD /sdh check-port --port 8080
COPY        --from=ghcr.io/sr-g/simple-docker-healthcheck:latest /sdh /sdh
```

In a Docker Compose : 

```
services:
    my_service:
        image: <my_image>
        healthcheck:
            test: ["CMD", "/sdh", "check-port", "--port", "8080"]
            interval: 60s
            timeout: 30s
            retries: 5
            start_period: 30s  
```

## DEV Activities

### Build

#### Build the binary

```bash
$ CGO_ENABLED=0 go build -ldflags "-s -w" .
$ upx -9 ./simple-docker-healthcheck
```

Or use the makefile : 

```bash
make build
```

#### Build the docker image

```bash
make docker-build
(...)
make docker-run
```

#### Release

```bash
TAG="v1.1.0-RELEASE"
git add .
git commit -m"chore(release): prepare release ${TAG}"
git push github master
git tag ${TAG}
git push github ${TAG}
```

### Overwrite a previous tag

```bash
TAG="v1.1.0-RELEASE"
git tag --delete ${TAG}
git push --delete github ${TAG}
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

- `docker-healthcheck` : a similar preliminary implementation and has been a source of inspiration (link : https://github.com/bratteng/docker-healthcheck) (the original project has been archived, this one is not a fork but a full/more complete rewrite)
- usage of that `docker-healthcheck` binary in the NGINX docker image, directly from the docker container, as an example of usage : [https://github.com/bratteng/docker-nginx/...](https://github.com/bratteng/docker-nginx/blob/15ddec93d6a47ca04f84cdf3bde8b834dee1b806/Dockerfile#L177-L178)
- General discussion on healthcheck with distroless containers : https://itnext.io/healthchecks-with-distroless-containers-262a52abc31e