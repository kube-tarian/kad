FROM golang:1.19 AS builder
WORKDIR /
COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/config_worker config-worker/main.go

FROM scratch
COPY --from=builder ./build/config_worker /config_worker
ENTRYPOINT ["/config_worker"]
