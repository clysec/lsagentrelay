FROM golang:1.26rc2-alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY lsagentrelay /app/lsagentrelay

ENTRYPOINT ["/app/lsagentrelay"]