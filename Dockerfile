FROM golang:1.16-buster AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o session-checker ./cmd/main.go
RUN useradd -u 10001 app


FROM scratch

COPY --from=builder /build/session-checker /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

USER app

CMD ["/session-checker"]
