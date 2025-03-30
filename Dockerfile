# Build stage
FROM --platform=linux/amd64 golang:1.21-alpine AS builder
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o applications_form_bot ./applications_form_bot
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o consultation_bot ./consultation_bot

# Run stage
FROM --platform=linux/amd64 alpine:latest
WORKDIR /app
COPY --from=builder /app/applications_form_bot/applications_form_bot .
COPY --from=builder /app/applications_form_bot/imgs ./imgs
COPY --from=builder /app/consultation_bot/consultation_bot .
COPY --from=builder /app/consultation_bot/imgs ./imgs
COPY --from=builder /app/.env .
COPY ./start.sh .

RUN chmod +x /app/applications_form_bot
RUN chmod +x /app/consultation_bot
RUN chmod +x /app/start.sh

CMD ["/app/start.sh"]