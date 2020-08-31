FROM golang:alpine AS builder
WORKDIR /workdir
COPY . .
RUN go build -ldflags '-w -s' -o /app/server .

FROM alpine
COPY --from=builder /app/server /usr/bin/server
EXPOSE 8080
ENTRYPOINT [ "server" ]
