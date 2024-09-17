FROM golang:1.23.0-alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY lsagentrelay /app/lsagentrelay

ENTRYPOINT ["/app/lsagentrelay"]