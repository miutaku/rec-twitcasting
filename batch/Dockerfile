# Build stage
FROM golang:1.23 AS builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY batch .
RUN go build -o batch-rec-twitcasting ./main.go

# Runtime stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/batch-rec-twitcasting .
RUN apk add --no-cache libc6-compat
CMD ["/root/batch-rec-twitcasting"]
