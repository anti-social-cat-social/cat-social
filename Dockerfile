################################################### STAGE 1
FROM golang:1.22-alpine AS builder

LABEL maintainer="Murz"

WORKDIR /app/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o main ./main.go

################################################### STAGE 2
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./
COPY --from=builder /app/.env ./

ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["./main"]