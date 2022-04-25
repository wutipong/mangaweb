#Stage 1
FROM golang:1.18-alpine AS builder

WORKDIR /go/src/mangaweb
COPY . .

ARG VERSION=Development
RUN apk add git
RUN go get -d -v ./...
RUN go build -v -ldflags="-X 'main.versionString=$VERSION' " -o mangaweb .

# Stage 2

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /go/src/mangaweb/mangaweb ./
COPY --from=builder /go/src/mangaweb/static ./static
COPY --from=builder /go/src/mangaweb/template ./template

CMD ["./mangaweb"]