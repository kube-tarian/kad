FROM cgr.dev/chainguard/go:latest AS builder
WORKDIR /
COPY ./capten ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/agent agent/cmd/main.go

FROM scratch
COPY --from=builder ./build/agent agent
COPY ./capten/database/cassandra/ cassandra/

USER 65532:65532
ENTRYPOINT ["./agent"]
