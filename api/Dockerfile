# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o api-rec-twitcasting ./main.go

# Runtime stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api-rec-twitcasting .
RUN apk add --no-cache libc6-compat ffmpeg
EXPOSE 8080
CMD ["/root/api-rec-twitcasting"]
