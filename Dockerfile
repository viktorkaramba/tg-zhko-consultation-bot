# Build stage
FROM --platform=linux/amd64 golang:1.21-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum to benefit from Docker cache layers
COPY go.mod go.sum ./
RUN go mod download

# Copy and build consultation bot
COPY ./consultation_bot ./consultation_bot
WORKDIR /app/consultation_bot
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o consultation_bot

# Copy and build applications form bot
WORKDIR /app
COPY ./applications_form_bot ./applications_form_bot
WORKDIR /app/applications_form_bot
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o applications_form_bot

# Run stage
FROM --platform=linux/amd64 alpine:latest
WORKDIR /app

# Install supervisor to manage multiple processes
RUN apk add --no-cache supervisor

# Copy configuration files
COPY supervisord.conf /etc/supervisord.conf
COPY .env ./.env

# Copy consultation bot files
COPY --from=builder /app/consultation_bot/consultation_bot ./consultation_bot
COPY ./consultation_bot/imgs ./consultation_bot/imgs

# Copy applications form bot files
COPY --from=builder /app/applications_form_bot/applications_form_bot ./applications_form_bot
COPY ./applications_form_bot/imgs ./applications_form_bot/imgs

# Set proper permissions
RUN chmod +x /app/consultation_bot
RUN chmod +x /app/applications_form_bot

# Set health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 CMD pgrep supervisord || exit 1

# Use supervisor to run both bots
CMD ["supervisord", "-c", "/etc/supervisord.conf"] 