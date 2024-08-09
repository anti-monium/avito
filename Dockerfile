FROM golang:1.22

COPY . /go/src/app

WORKDIR /go/src/app/cmd

RUN go build -o avito main.go

EXPOSE 9090

CMD ["./avito"]