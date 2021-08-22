FROM golang:1.17-alpine AS builder
WORKDIR /workdir
COPY . .
RUN go build -ldflags '-w -s' -o /app/echoserver .

FROM alpine:3.14.1
COPY --from=builder /app/echoserver /usr/bin/echoserver
EXPOSE 8080
ENTRYPOINT [ "echoserver" ]
