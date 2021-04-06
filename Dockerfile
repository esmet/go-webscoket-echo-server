FROM golang:latest as builder
WORKDIR /app
COPY frontend.go go.mod go.sum main.go ./

RUN go build -o /go-websocket-echo-server ./

FROM alpine:latest as build
RUN apk add --no-cache ca-certificates

FROM docker.io/frolvlad/alpine-glibc:alpine-3.12_glibc-2.32
COPY --from=builder /go-websocket-echo-server /go-websocket-echo-server
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
ENV PORT 8080
EXPOSE 8080
USER nobody
ENTRYPOINT ["/go-websocket-echo-server"]
