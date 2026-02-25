FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
USER appuser
WORKDIR /app
COPY --from=build /bin/api /app/api
EXPOSE 8080
CMD ["/app/api"]
