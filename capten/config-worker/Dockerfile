FROM golang:1.19 AS builder
WORKDIR /
COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/config_worker config-worker/main.go

FROM alpine:3.16
RUN apk add --no-cache libc6-compat
COPY --from=builder ./build/config_worker config_worker
COPY config-worker/git-config-files /git-config-files

ENTRYPOINT ["./config_worker"]
