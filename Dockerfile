FROM golang:1.17-alpine

WORKDIR /go/src/mangaweb
COPY . .

ARG VERSION=Development
RUN go get -d -v ./...
RUN go install -v -ldflags="-X 'main.versionString=$VERSION' " ./... 

CMD ["mangaweb"]