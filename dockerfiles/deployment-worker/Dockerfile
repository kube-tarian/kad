FROM golang:1.19.4 AS builder
WORKDIR /
COPY ./integrator ./
RUN rm -rf vendor

RUN go mod download
RUN go build -o ./build/deployment-worker deployment-worker/main.go

FROM alpine:3.16
RUN apk add --no-cache libc6-compat
COPY --from=builder ./build/deployment-worker deployment-worker

USER 65532:65532
ENTRYPOINT ["./deployment-worker"]
