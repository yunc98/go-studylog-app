# 一旦最新バージョンを使用
FROM golang:latest

# アップデートとgitのインストール
RUN apt-get update && apt-get install -y git

WORKDIR /go/src

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

CMD ["go", "run", "main.go"]