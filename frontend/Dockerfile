# -----------------------------
# Stage 1 — Build frontend
# -----------------------------
FROM node:20-bookworm AS builder

WORKDIR /build

# Copy only what we need for building
COPY package.json package-lock.json* ./

# Install dependencies first (better caching)
RUN npm ci --ignore-scripts || npm install

# Copy protobuf files and source code together
COPY gen/ts ./gen/ts
COPY . .

# Debug: list directory contents
RUN ls -la /build && echo "---" && ls -la /build/gen/ts 2>/dev/null || echo "No gen/ts found"

# Build static site with error output
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN npm run build 2>&1 || (echo "Build failed, checking package.json..." && cat package.json)


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