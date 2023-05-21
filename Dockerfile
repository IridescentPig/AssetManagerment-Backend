FROM golang:1.20

ENV GIN_MODE release

ENV RELEASE 1

WORKDIR /opt/tmp

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build .

EXPOSE 80
EXPOSE 8080

CMD ["./asset-management"]