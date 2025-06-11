FROM golang:1.24.3

WORKDIR /app

RUN go install github.com/air-verse/air@v1.62.0 \
    && go install golang.org/x/tools/cmd/goimports@v0.34.0 \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
       | sh -s -- -b /go/bin v2.1.6 \
    && go install github.com/swaggo/swag/cmd/swag@v1.16.4 \
    && apt-get update && apt-get install -y make && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

CMD ["air"]
