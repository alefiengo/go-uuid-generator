FROM golang:alpine3.19 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM scratch

COPY --from=build /app/main /app/main

CMD ["/app/main"]
