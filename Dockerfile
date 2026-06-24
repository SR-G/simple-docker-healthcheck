# FROM        scratch
FROM        gcr.io/distroless/static-debian12

ENTRYPOINT  ["/sdh"]

ADD        ./sdh /sdh