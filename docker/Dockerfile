FROM golang:1.24-alpine AS builder

LABEL maintainer="jasonweiss@gmail.com"
LABEL description="Build stage for Go application"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o myapp .

FROM alpine:latest

LABEL description="Production image for Go application"

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app/

COPY --from=builder /app/myapp .

EXPOSE 50051

ENTRYPOINT ["./myapp", "server"]