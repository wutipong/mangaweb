    FROM golang:1.14
    WORKDIR /app
    COPY . .
    RUN go build
    CMD ["/app/mangaweb"]