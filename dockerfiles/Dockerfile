# Use the latest Alpine as the runtime layer, pinned to the build platform
FROM --platform=$BUILDPLATFORM alpine:latest AS alpine

# Install certificates and timezone data so TLS and timezones work at runtime
RUN apk add --no-cache \
    ca-certificates \
    tzdata

FROM scratch
LABEL "com.centurylinklabs.watchtower"="true"

COPY --from=alpine \
    /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/ca-certificates.crt
COPY --from=alpine \
    /usr/share/zoneinfo \
    /usr/share/zoneinfo

COPY watchtower /

HEALTHCHECK CMD [ "/watchtower", "--health-check"]

ENTRYPOINT ["/watchtower"]
