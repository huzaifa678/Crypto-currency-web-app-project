FROM golang:1.25.8-alpine3.23 AS builder
WORKDIR /app
COPY . .

# static build for minimal final image
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main main.go

# Match the builder's Alpine version and stay on a supported release.
FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./db/migrations

# Pull the latest security patches, then create an unprivileged user.
RUN apk --no-cache upgrade \
  && addgroup -S app && adduser -S app -G app \
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