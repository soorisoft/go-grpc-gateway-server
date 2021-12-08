FROM golang:1.16.3-alpine3.13

ENV GO111MODULE=on

RUN mkdir /app
## Copy everything in the root directory into our /app directory
ADD . /app

## We specify that we now wish to execute
## any further commands inside our /app
## directory
WORKDIR /app

RUN go mod download

RUN go build -o main .

EXPOSE 9091
## Our start command which kicks off
## our newly created binary executable
CMD ["/app/main"]
