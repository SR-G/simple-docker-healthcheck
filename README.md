

## Cookbooks

```bash
go run . check-port --hostname 192.168.8.140 --port 32400
```

```bash
go run . check-http-code --url "https://gluetun-torrents-exporter.tensin.ovh/metrics" 
--min-status-code 0 --max-status-code 399
```

```bash
go run . check-http-code --url "https://gluetun-torrents-exporter.tensin.ovh/metrics" --status-code 200
```

```bash
go run . check-http-text --url "https://gluetun-torrents-exporter.tensin.ovh/metrics" --text "go_gc_duration_seconds_sum"
```

```bash
go run . check-http-json --url "https://time.tensin.org/" --json-path '$.date' --value '2026-05-16'
```

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


## Build

```bash
$ go build -ldflags "-s -w" .
$ upx -9 ./simple-docker-healthcheck
```

| Binary                                             | Ratio | Size    |
| -------------------------------------------------- | ----- | ------- |
| Built with go compiler                             | 100%  | 13.3 M  |
| Built with go compiler + upx -9                    |  56%  |  7.0 M  |
| Built with go compiler with ldflags -s -w          |       |  8.3 M  |
| Built with go compiler with ldflags -s -w + upx -9 |       |  3.3 M  |


```bash
CGO_ENABLED=0 docker build . -t simple-docker-healthcheck
(...)
docker run -it --rm --name "test" simple-docker-healthcheck --version
```

```
IMAGE                              ID             DISK USAGE   CONTENT SIZE   EXTRA
node:24-alpine                     edd927012c1e        160MB             0B        
simple-docker-healthcheck:latest   96ea9531359e       5.46MB             0B        
```