FROM golang:1.19.4 AS builder
WORKDIR /
COPY ./server ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/server cmd/server/main.go

FROM alpine:3.16
RUN apk add --no-cache libc6-compat
COPY --from=builder ./build/server server
COPY --from=builder ./server/storeconfig storeconfig

USER 65532:65532
ENTRYPOINT ["./server"]
