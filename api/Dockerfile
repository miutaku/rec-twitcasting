# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go build -o rec-twitcasting ./api/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/rec-twitcasting .
RUN apk add --no-cache libc6-compat ffmpeg
EXPOSE 8080
CMD ["/root/rec-twitcasting"]
