FROM golang:1.24.3 AS base

WORKDIR /app

RUN go install github.com/air-verse/air@v1.62.0 \
    && go install golang.org/x/tools/cmd/goimports@latest \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
       | sh -s -- -b /go/bin v1.64.8 \
    && go install github.com/swaggo/swag/cmd/swag@v1.16.4

ENV PATH=$PATH:/go/bin

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
COPY .air.toml ./

CMD ["air"]
