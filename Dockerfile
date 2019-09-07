FROM golang:1.13.0-alpine3.10 AS builder
WORKDIR /src
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags '-w -s' -o /app/k8s-example-server .

FROM alpine:3.10 as production
COPY --from=builder /app/k8s-example-server /usr/bin/k8s-example-server
ENTRYPOINT [ "/usr/bin/k8s-example-server" ]
