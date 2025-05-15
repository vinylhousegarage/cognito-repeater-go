FROM golang:1.24
WORKDIR /app
COPY . .
RUN go install github.com/cosmtrek/air@latest
CMD ["air"]
