FROM golang:1.24.3-alpine3.21 AS builder
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

RUN chmod +x start.sh wait-for.sh

EXPOSE 8081
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]