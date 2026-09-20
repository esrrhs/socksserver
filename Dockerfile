# Build stage
FROM golang:alpine AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o socksserver .

# Runtime stage
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /
COPY --from=builder /build/socksserver /socksserver
COPY --from=builder /build/socksserver /usr/local/bin/socksserver
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/socksserver", "-h"]
