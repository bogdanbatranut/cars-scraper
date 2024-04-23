# Start from the latest golang base image
FROM golang:alpine as builder

# Install Essentials
RUN apk update \
    && apk add -U --no-cache ca-certificates \
    && update-ca-certificates
# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scraper cmd/newversion/scraper.go

RUN ls -lah
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/scraper .

ENV APP_BASE_URL=dev.auto-mall.ro APP_DB_USER=dev APP_DB_PASS=siana1316 APP_DB_NAME=automall APP_DB_HOST=dev.auto-mall.ro SMQ_URL=dev.auto-mall.ro BACKEND_HTTP_PORT=8080

ENTRYPOINT ["./scraper"]