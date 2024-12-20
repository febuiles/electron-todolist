FROM golang:1.21-alpine

WORKDIR /app
COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN rm todos.db
RUN apk add --no-cache gcc musl-dev sqlite-dev

RUN go build -o server .

EXPOSE 8080

CMD ["./server"]
