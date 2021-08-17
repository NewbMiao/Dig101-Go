
FROM golang:1.17-alpine
RUN apk update && apk add --update git 
RUN go get -u golang.org/x/perf/cmd/benchstat 
RUN go get -u github.com/NewbMiao/Dig101-Go
CMD GOGC=off GODEBUG=asyncpreemptoff=1 go test -gcflags='-N -l' github.com/NewbMiao/Dig101-Go/types/struct -bench . -count 20 -cpu 1 > b.txt && benchstat b.txt

# This benchmark is useless now. Needs figure out a new way to do it.
# docker build -t  gobench-structalign .
# docker run --rm   gobench-structalign