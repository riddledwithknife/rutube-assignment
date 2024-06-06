FROM golang:alpine AS builder

WORKDIR /birthday-greetings

COPY . ./

RUN go mod tidy

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /birthday-greetings/main .

EXPOSE 8080

CMD ["./main"]
