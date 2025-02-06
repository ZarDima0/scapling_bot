FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o scalping-bot ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/scalping-bot .
COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./scalping-bot"]
