FROM golang:1.20

ENV GIN_MODE release

WORKDIR /opt/tmp

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

EXPOSE 8080

CMD ["go", "run", "main.go"]