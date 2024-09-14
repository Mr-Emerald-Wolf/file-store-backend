# Build stage
FROM golang:1.21.0-alpine3.18 AS build
WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -o main ./cmd

# Run stage  
FROM alpine:3.18
RUN apk update --no-cache && apk add --no-cache ca-certificates
WORKDIR /app
COPY .env .env
COPY --from=build /app/main .
EXPOSE 8080
CMD ["./main"]
