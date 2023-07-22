#Stage 1 -- building executable
FROM golang:1.20-alpine AS builder1

WORKDIR /go/src/mangaweb
COPY . .

ARG VERSION=Development
RUN apk add git
RUN go get -d -v ./...
RUN go build -v -ldflags="-X 'main.versionString=$VERSION' " -o mangaweb .

# Stage 2 -- building resources
FROM node:20-alpine AS builder2

WORKDIR /go/src/mangaweb
COPY . .

RUN npm install
RUN npm run build

# Stage 3 -- combine the two
FROM alpine:latest

WORKDIR /root/
COPY --from=builder1 /go/src/mangaweb/mangaweb ./
COPY --from=builder2 /go/src/mangaweb/static ./static
COPY --from=builder1 /go/src/mangaweb/template ./template

CMD ["./mangaweb"]