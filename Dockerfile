# Use bullseye base image as it has upx available
FROM golang:1.24-bullseye AS builder

LABEL maintainer="github.com/shiftavenue"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apt-get -qq update \
    && apt-get -yqq install upx

COPY . .

RUN go build -ldflags "-s -w" -o /bin/gas-action ./cmd/gas-action \
    && upx -q -9 /bin/gas-action

RUN echo "nobody:topsecret:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc_passwd /etc/passwd
COPY --from=builder --chown=65534:0 /bin/gas-action /gas-action

USER nobody
ENTRYPOINT ["/gas-action"]
