FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o app main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/app .

CMD ["./app"]
