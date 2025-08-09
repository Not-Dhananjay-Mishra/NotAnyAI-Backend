FROM golang:1.24.5
WORKDIR /app
COPY . .
RUN go build -v ./...
CMD ["./server"]
