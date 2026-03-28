FROM golang:1.26.1

WORKDIR /app

COPY go.mod ./

RUN go mod tidy

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8080

CMD [ "/docker-gs-ping" ]