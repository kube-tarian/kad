FROM golang:1.19 AS builder
# All these steps will be cached
RUN mkdir /agent
WORKDIR /agent

# COPY the source code as the last step
COPY . .
RUN go mod vendor
RUN go build -o intelopsagent cmd/agent/main.go

FROM alpine:3.16
RUN apk add --no-cache libc6-compat
RUN mkdir /agent
WORKDIR /agent
COPY --from=builder /agent/intelopsagent .
ENTRYPOINT ["./intelopsagent"]