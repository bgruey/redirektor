FROM golang:1.22

WORKDIR /server

COPY ./ ./

RUN go mod tidy

RUN go build -o /app.x ./server.go

CMD ["/app.x"]