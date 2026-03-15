# -----------------------------
# Stage 1 — Build frontend
# -----------------------------
FROM node:20-bookworm AS builder

WORKDIR /build

# Copy package files first for caching
COPY package.json package-lock.json* ./
RUN npm ci || npm install

# Copy all source including gen/ts (it's in the build context)
COPY . .

# Build static site with verbose output
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN npm run build 2>&1 | tee /tmp/build.log || (echo "Build failed, logs:" && cat /tmp/build.log && exit 1)


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