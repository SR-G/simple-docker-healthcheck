# FROM        scratch
FROM        gcr.io/distroless/static-debian12

ENTRYPOINT  ["/simple-docker-healthcheck"]

ADD        ./simple-docker-healthcheck /