FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./db/migrations

# Create an unprivileged user and run the container as it.
RUN addgroup -S app && adduser -S app -G app \
  && chmod +x start.sh wait-for.sh \
  && chown -R app:app /app

USER app

EXPOSE 8081
EXPOSE 9090

# Report container health by probing the HTTP server port.
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD wget -q --spider http://localhost:8081/ || exit 1

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]