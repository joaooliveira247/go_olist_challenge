FROM golang:1.23-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY ./
RUN go mod download

RUN go run main.go db create

CMD ["air", "run"]