FROM golang:1.19.4 AS builder
WORKDIR /
COPY ./capten ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/config-worker config-worker/main.go

FROM scratch
COPY --from=builder ./build/config-worker /config-worker

USER 65532:65532
ENTRYPOINT ["./config-worker"]
