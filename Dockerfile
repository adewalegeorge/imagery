# Build stage
FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y pkg-config libvips-dev
WORKDIR /app/api
RUN go build -o /app/api/app .

# Runtime stage
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y libvips-dev && rm -rf /var/lib/apt/lists/*
WORKDIR /app/api
COPY --from=builder /app/api/app ./app
EXPOSE 8080
CMD ["./app"] 