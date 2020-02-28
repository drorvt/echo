FROM golang:1.14-alpine AS builder
WORKDIR /src
COPY . .
RUN go build -ldflags '-w -s' -o /app/server .

FROM alpine:3.11
COPY --from=builder /app/k8s-example-server /usr/bin/server
EXPOSE 10200
ENTRYPOINT [ "/usr/bin/server" ]
