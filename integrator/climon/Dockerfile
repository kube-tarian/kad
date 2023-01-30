FROM golang:1.19 AS builder
WORKDIR /
COPY ./ ./

RUN go mod download
RUN CGO_ENABLED=0 go build -o ./build/climon climon/main.go

FROM scratch
COPY --from=builder ./build/climon climon

USER 65532:65532
ENTRYPOINT ["./climon"]
