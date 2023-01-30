FROM golang:1.19.4 AS builder
WORKDIR /
COPY ./integrator ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/deployment-worker deployment-worker/main.go

FROM scratch
COPY --from=builder ./build/deployment-worker deployment-worker

USER 65532:65532
ENTRYPOINT ["./deployment-worker"]
