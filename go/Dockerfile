# バージョン指定
FROM golang:1.16

WORKDIR /go/src

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -v -o app

CMD ["./app"]