# Build stage
FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y pkg-config libvips-dev fonts-dejavu-core
RUN go build -o /app/api/app ./api

# Runtime stage
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y libvips-dev fonts-dejavu-core && rm -rf /var/lib/apt/lists/*
WORKDIR /app/api
COPY --from=builder /app/api/app ./app
EXPOSE 8080
CMD ["./app"]