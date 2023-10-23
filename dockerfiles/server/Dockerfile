FROM golang:1.20 AS builder
WORKDIR /
COPY ./server ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/server cmd/server/main.go

FROM alpine:3.16
RUN apk add --no-cache libc6-compat
COPY --from=builder ./data /data
COPY --from=builder ./build/server server

USER 65532:65532
ENTRYPOINT ["./server"]
