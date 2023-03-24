FROM golang:1.20-alpine

WORKDIR /pay-with-crypto

COPY . ./

RUN apk update
RUN apk add postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o main .
EXPOSE 8081

CMD ["./main", "-p"]