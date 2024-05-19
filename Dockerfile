FROM golang:1.22
# Set destination for COPY
WORKDIR /app
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app

FROM alpine:latest
RUN apk add -u ca-certificates
RUN adduser -D -h /app -u 1000 -k /dev/null app
WORKDIR /app
# Copy the binary from the builder image
COPY --from=0 /app/bin/app bin/app
# Run
CMD ["bin/app"]
