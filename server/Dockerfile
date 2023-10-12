FROM golang:1.20
WORKDIR /app
COPY . .

RUN go build -o server server/cmd

EXPOSE 5000
CMD ["./server"]