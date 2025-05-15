FROM golang:1.24.3 AS base
FROM cosmtrek/air
WORKDIR /app
COPY . .
CMD ["air"]
