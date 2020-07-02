FROM golang

ADD ./wait_for_it.sh /
ADD . /go/src/contree

WORKDIR /go/src/contree

RUN go install ./...

RUN go build

ENTRYPOINT /wait_for_it.sh