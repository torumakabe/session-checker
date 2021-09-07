FROM golang:1.17-bullseye AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make build
RUN useradd -u 10001 app


FROM scratch

COPY --from=builder /build/bin/session-checker /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

ENV GIN_MODE=release
USER app

CMD ["/session-checker"]
