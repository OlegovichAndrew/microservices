# syntax=docker/dockerfile:1
FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/scooter_service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/scooter_service ./cmd

FROM scratch
COPY --from=builder /go/src/scooter_service/build/scooter_service ./usr/bin/scooter_service
ENTRYPOINT [ "/usr/bin/scooter_service" ]