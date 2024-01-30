FROM golang:latest

LABEL authors="Alejandro Fiengo"

WORKDIR /app

COPY . .

RUN go build -o main .

CMD ["./main"]
