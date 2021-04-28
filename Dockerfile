FROM golang:1.16-alpine

WORKDIR /go/src/mangaweb
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["mangaweb"]