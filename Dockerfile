# STAGE 1: Build
FROM golang:alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Compile with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o rss-reader cmd/server/main.go

# STAGE 2: Production Image
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/rss-reader .

# COPY THE MISSING FILES!
# We need the SQL schema for the DB, and the assets for the Tailwind CSS
COPY --from=builder /app/db/schema.sql ./db/schema.sql
COPY --from=builder /app/ui/assets ./ui/assets

# Create the directory for the SQLite volume
RUN mkdir -p /app/data

# Tell the app it is running in production
ENV ENV=production

EXPOSE 8080
CMD ["./rss-reader"]