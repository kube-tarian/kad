FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /
COPY ./server ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/server cmd/server/main.go

FROM cgr.dev/chainguard/wolfi-base
RUN apk add --no-cache glibc
COPY --from=builder ./data /data
COPY --from=builder ./build/server server

USER 65532:65532
ENTRYPOINT ["./server"]
