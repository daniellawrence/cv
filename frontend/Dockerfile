# -----------------------------
# Stage 1 — Build frontend
# -----------------------------
FROM node:20-bookworm AS builder

WORKDIR /build

# Copy package files first for caching
COPY package.json package-lock.json* ./

# Install dependencies, but skip the file: dependency for now
RUN npm ci --ignore-scripts || npm install --no-optional

# Manually copy and link @cv/proto
RUN mkdir -p node_modules/@cv && cp -r ../gen/ts node_modules/@cv/proto 2>/dev/null || echo "Proto files may not be available yet"

# Copy source including gen/ts (in build context)
COPY . .

# Reinstall with all dependencies now that proto is in place
RUN npm install

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
COPY nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]