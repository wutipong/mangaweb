FROM golang:1.17-alpine

WORKDIR /go/src/mangaweb
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["mangaweb"]