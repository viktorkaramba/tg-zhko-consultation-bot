# Build stage
FROM --platform=linux/amd64 golang:1.21-alpine AS builder
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bot ./consultation_bot

# Run stage
FROM --platform=linux/amd64 alpine:latest
WORKDIR /app
COPY --from=builder /app/bot .
COPY --from=builder /app/consultation_bot/imgs ./imgs
COPY --from=builder /app/.env .

RUN chmod +x /app/bot
CMD ["/app/bot"] 