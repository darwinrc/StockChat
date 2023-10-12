FROM golang:1.20
WORKDIR /app
COPY . .

RUN go build -o bot bot/cmd

EXPOSE 5000
CMD ["./bot"]