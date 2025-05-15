FROM golang:1.24.3 AS base
WORKDIR /app
COPY . .
RUN test -f go.mod || go mod init cognito-repeater-go
RUN go mod tidy
RUN go install github.com/air-verse/air@latest
ENV PATH=$PATH:/go/bin
CMD ["air"]
