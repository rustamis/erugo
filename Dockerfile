FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ .
RUN npm run build

FROM node:18-alpine AS docs-builder
WORKDIR /app
COPY erugo.openapi.json .
RUN npx -y @redocly/cli build-docs erugo.openapi.json --output=swagger.html

FROM golang:1.21-alpine AS backend-builder

WORKDIR /app
COPY . .
COPY --from=frontend-builder /app/frontend/dist frontend/dist
COPY --from=docs-builder /app/swagger.html .

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o erugo

FROM alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /app
COPY --from=backend-builder /app/erugo .
COPY --from=backend-builder /app/config.json .

# Create directories for storage and database
RUN mkdir -p /app/storage /app/private

# Set environment variables with defaults
ENV ERUGO_BASE_STORAGE_PATH=/app/storage \
    ERUGO_APP_URL=http://localhost:9199 \
    ERUGO_BIND_PORT=9199 \
    ERUGO_DATABASE_FILE_PATH=/app/private/erugo.db \
    ERUGO_MAX_SHARE_SIZE=2G \
    ERUGO_PRIVATE_DATA_PATH=/app/private

EXPOSE 9199

CMD ["./erugo"]