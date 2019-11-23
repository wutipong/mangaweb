    FROM golang:1.13
    WORKDIR /app
    COPY . .
    RUN go build
    CMD ["/app/mangaweb"]