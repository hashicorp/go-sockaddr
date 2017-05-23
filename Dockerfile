FROM golang:alpine

ADD ./ $GOPATH/go-sockaddr

RUN set -x && \
    apk add --update git && \
    go get -v -u github.com/hashicorp/go-sockaddr/cmd/sockaddr

CMD [ "sockaddr" ]
