FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o authentication-service .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/authentication-service .

EXPOSE 8088

CMD ["./authentication-service"]
