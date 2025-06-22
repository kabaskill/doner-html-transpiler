# Build stage for frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Clean install with explicit dev dependencies
RUN npm ci --include=dev

# Copy source code
COPY frontend/ ./

# Build with vite only (vite handles TypeScript compilation)
RUN npx vite build

# Build stage for backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/backend

# Copy go mod files first for better caching
COPY backend/go.mod ./
COPY backend/go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy backend binary
COPY --from=backend-builder /app/backend/main .

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./static

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
