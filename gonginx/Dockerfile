# build stage
FROM golang:1.21.5-alpine3.19 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go

# run stage
FROM alpine:3.19 AS runner
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080