# Build stage for frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ ./
RUN npm run build

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
