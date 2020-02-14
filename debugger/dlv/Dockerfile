FROM golang:1.14rc1
WORKDIR /go/src/app
RUN go get -u github.com/go-delve/delve/cmd/dlv
COPY . .
CMD ["app"]

