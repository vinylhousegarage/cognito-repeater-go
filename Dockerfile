FROM golang:1.24.3 AS base

WORKDIR /app

RUN go install github.com/air-verse/air@latest \
    && go install golang.org/x/tools/cmd/goimports@latest \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
       | sh -s -- -b /go/bin v2.1.6

ENV PATH=$PATH:/go/bin

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

CMD ["air"]
