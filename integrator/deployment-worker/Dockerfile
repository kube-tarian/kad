FROM golang:1.19 AS builder
WORKDIR /
COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/deployment_worker deployment-worker/main.go

FROM scratch
COPY --from=builder ./build/deployment_worker deployment_worker

USER 65532:65532
ENTRYPOINT ["./deployment_worker"]
