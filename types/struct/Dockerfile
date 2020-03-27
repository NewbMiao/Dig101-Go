
FROM golang:1.14-alpine
RUN apk update && \
    apk add --update git 
RUN go get -u golang.org/x/perf/cmd/benchstat && go get -u github.com/aclements/perflock/cmd/perflock && install $GOPATH/bin/perflock /usr/bin/perflock
RUN go get -u github.com/NewbMiao/Dig101-Go
CMD (perflock -daemon &) && sleep 1 \
  && perflock go test -gcflags='-N -l' github.com/NewbMiao/Dig101-Go/types/struct -bench . -count 3 > old.txt \
  && benchstat old.txt


# docker build -t  gobench-structalign .
# docker run --rm   gobench-structalign