FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o application-service .

FROM alpine:latest

RUN apk --no-cache add sqlite

WORKDIR /root/

COPY --from=builder /app/application-service .
COPY --from=builder /app/todos.db .

EXPOSE 8087

CMD ["./application-service"]
