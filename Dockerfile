FROM golang:1.13.6-alpine as builder

# ENV GO111MODULE=on

RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /payments

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src
COPY .env .

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main ./src

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /payments/main .
COPY --from=builder /payments/.env .

EXPOSE 8095
CMD ["./main"]