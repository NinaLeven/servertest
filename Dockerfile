FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app

COPY . .

RUN go build -o servertest servertest/bin

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/ /app/

CMD /app/servertest -c /app/config/dev.yaml