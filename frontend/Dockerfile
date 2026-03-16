# Build stage - compile the React application
FROM node:20-alpine AS builder

WORKDIR /app

# Copy all files (including gen directory) for local dependency resolution
COPY . .

# Install frontend dependencies (includes @cv/proto local dependency)
RUN cd frontend && npm ci

# Build the frontend application
WORKDIR /app/frontend
RUN npm run build

# Production stage - serve with nginx
FROM nginx:alpine AS production

# Install curl for health checks and envsubst
RUN apk add --no-cache curl gettext

# Copy custom nginx configuration template
COPY frontend/nginx.conf.template /etc/nginx/conf.d/default.conf.template

# Copy docker entrypoint script
COPY frontend/docker-entrypoint.sh /docker-entrypoint.d/99-envsubst-nginx.sh
RUN chmod +x /docker-entrypoint.d/99-envsubst-nginx.sh

# Copy built files from the builder stage
COPY --from=builder /app/frontend/dist /usr/share/nginx/html

# Expose port 8080
EXPOSE 8080

# Environment variables will be set at runtime via docker run -e or in compose files
# For build-time defaults, use ARG with default values
ARG BUILD_TIMESTAMP=unknown
ARG CONTENT_SHA1=none

ENV BUILD_TIMESTAMP=$BUILD_TIMESTAMP
ENV CONTENT_SHA1=$CONTENT_SHA1

# Health check endpoint (optional)
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
