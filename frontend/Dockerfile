# -----------------------------
# Stage 1 — Build frontend
# -----------------------------
FROM node:20-alpine AS builder

WORKDIR /build

# Copy protobuf generated files first (required for build)
COPY gen/ts ./gen/ts

# Install dependencies first (better caching)
COPY package.json package-lock.json* ./
RUN npm ci --ignore-scripts 2>/dev/null || npm install

# Copy source
COPY . .

# Build static site
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN npm run build


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
COPY frontend/nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]