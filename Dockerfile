FROM golang:1.24.3 AS base
WORKDIR /app
RUN go install github.com/air-verse/air@latest \
    && go install golang.org/x/tools/cmd/goimports@latest
ENV PATH=$PATH:/go/bin
CMD ["air"]
