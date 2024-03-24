# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o main .

FROM seleniarm/standalone-chromium
WORKDIR /app
COPY --from=builder /app/main .
USER 1001
CMD [ "/app/main" ]