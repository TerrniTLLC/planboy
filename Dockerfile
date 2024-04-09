FROM golang:latest

RUN mkdir /app

ADD . /app/
WORKDIR /app

# Get all dependencies without go mod download
# RUN go get github.com/crocone/tg-bot v0.0.0-20230406135341-8d3e6a822eaf
# RUN go get github.com/joho/godotenv v1.5.1
RUN go mod download


RUN go build -o main .

CMD ["/app/main"]

