# build stage
FROM golang:1.22-alpine3.19 AS builder

WORKDIR /app
COPY account .
COPY go.mod .
RUN go mod tidy
RUN go build -o main main.go

# run stage
FROM alpine:3.19 AS runner
WORKDIR /app
COPY --from=builder /app/main .
