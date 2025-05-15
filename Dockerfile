FROM golang:1.24.3 AS base
WORKDIR /app
COPY . .
RUN [ -f go.mod ] || go mod init cognito-repeater-go
RUN go mod tidy
FROM cosmtrek/air
WORKDIR /app
COPY --from=base /app /app
