FROM golang:1.14

RUN go get -u github.com/jackc/tern

WORKDIR /go/bin

COPY . .

CMD ["tern-migrate.sh"]
