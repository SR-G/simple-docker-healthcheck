

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

