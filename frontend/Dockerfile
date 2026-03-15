# -----------------------------
# Stage 1 — Build frontend
# -----------------------------
FROM node:20-bookworm AS builder

WORKDIR /build

# Copy package files first for caching
COPY package.json package-lock.json* ./

# Install all dependencies including file: dependency
RUN npm install

# Copy source and proto files (gen/ts is in build context at root)
COPY . .

# Build static site with error handling
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN set -e && echo "Starting build..." && npm run build && echo "Build completed successfully"


# -----------------------------
# Stage 2 — NGINX
# -----------------------------
FROM nginx:1.27-alpine

ARG CONTENT_SHA1=unknown
ARG BUILD_TIMESTAMP=unknown

ENV CONTENT_SHA1=${CONTENT_SHA1}
ENV BUILD_TIMESTAMP=${BUILD_TIMESTAMP}

# Remove default nginx content
RUN rm -rf /usr/share/nginx/html/*

# Copy built frontend
COPY --from=builder /build/dist /usr/share/nginx/html

# Custom nginx config (processed as envsubst template at startup)
COPY nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]