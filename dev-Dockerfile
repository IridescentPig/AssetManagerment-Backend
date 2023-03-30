FROM golang:1.20

ENV GIN_MODE debug

ENV DEBUG 1

WORKDIR /opt/tmp

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

EXPOSE 80

CMD ["go", "run", "main.go"]