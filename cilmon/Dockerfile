FROM golang:1.19 AS builder
# All these steps will be cached
RUN mkdir /cilmon
WORKDIR /cilmon

# COPY the source code as the last step
COPY . .
RUN go mod vendor
RUN go build -o worker cmd/climon.go

FROM alpine:3.16
RUN apk add --no-cache libc6-compat
RUN mkdir /cilmon
WORKDIR /cilmon
COPY --from=builder /cilmon/worker .
ENTRYPOINT ["./worker"]
