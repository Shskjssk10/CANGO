FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Use a minimal base image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

# Expose the port your application runs on
EXPOSE 8001

# Command to run the executable
CMD ["./main"]
