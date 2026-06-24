# FROM        scratch
FROM        gcr.io/distroless/static-debian12

ENTRYPOINT  ["/sdh"]

ADD        ./distribution/sdh-linux-amd64 /sdh