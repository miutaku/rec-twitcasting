# Build stage
FROM golang:1.23 AS builder
ENV CGO_ENABLED=0
WORKDIR /app
COPY manage-backend .
RUN go build -o manage-backend-rec-twitcasting ./main.go

# Runtime stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/manage-backend-rec-twitcasting .
RUN apk add --no-cache libc6-compat
CMD ["/root/manage-backend-rec-twitcasting"]
