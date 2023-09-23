FROM golang:1.19.4 AS builder
WORKDIR /
COPY ./capten ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/config-worker config-worker/main.go

FROM alpine:3.16
RUN apk add --no-cache libc6-compat
COPY --from=builder ./build/config-worker /config-worker
COPY --from=builder ./config-worker/git-config-files /git-config-files

ENTRYPOINT ["./config-worker"]
