FROM golang:1.24.3 AS base
WORKDIR /app
RUN go install github.com/air-verse/air@latest
ENV PATH=$PATH:/go/bin
CMD ["air"]
