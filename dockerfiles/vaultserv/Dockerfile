FROM golang:1.19.4 AS builder
RUN mkdir /vaultserv
WORKDIR /vaultserv
COPY ./vaultserv ./
RUN rm -rf vendor

RUN go mod download
RUN CGO_ENABLED=0 go build -o vaultserv cmd/main.go

FROM alpine:3.16
RUN apk add --no-cache libc6-compat
COPY --from=builder ./vaultserv/vaultserv vaultserv

USER 65532:65532
ENTRYPOINT ["./vaultserv"]