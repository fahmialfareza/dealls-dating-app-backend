# === Stage 1: Build ===
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY docker/production/.env.production .env

RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y libsodium-dev --no-install-recommends \
  && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o dating-app-server ./cmd/main.go

# === Stage 2: Packaging ===
FROM golang:1.23

# Set the time zone
ENV TZ=Asia/Jakarta

WORKDIR /app

# Install timezone data package and other dependencies
RUN apt-get update && apt-get install -y build-essential procps curl file git tzdata libsodium-dev locales --no-install-recommends \
  && apt-get clean && rm -rf /var/lib/apt/lists/*

# Create a non-root user and group
RUN groupadd -r nonroot && useradd -r -g nonroot -u 777 nonroot

# Copy the built application
COPY --from=builder /app/.env .
COPY --from=builder /app/dating-app-server .

# Change ownership of the application files to the non-root user
RUN chown -R nonroot:nonroot /app

# Switch to non-root user
USER nonroot

# Run the API
ENTRYPOINT ["./dating-app-server"]