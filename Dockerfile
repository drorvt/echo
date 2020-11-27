FROM golang:alpine AS builder
WORKDIR /workdir
COPY . .
RUN go build -ldflags '-w -s' -o /app/echoserver .

FROM alpine:latest
COPY --from=builder /app/echoserver /usr/bin/echoserver
EXPOSE 8080
ENTRYPOINT [ "echoserver" ]
