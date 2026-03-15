FROM node:18-alpine AS builder

WORKDIR /app/frontend

# Copy package files first for caching
COPY frontend/package.json frontend/package-lock.json* ./

# Install all dependencies including file: dependency
RUN npm install

# Copy source code
COPY . .

# Build static site
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN npm run build


FROM nginx:1.27-alpine

ARG CONTENT_SHA1=unknown
ARG BUILD_TIMESTAMP=unknown

ENV CONTENT_SHA1=${CONTENT_SHA1}
ENV BUILD_TIMESTAMP=${BUILD_TIMESTAMP}

# Remove default nginx content
RUN rm -rf /usr/share/nginx/html/*

# Copy built frontend
COPY --from=builder /app/frontend/dist /usr/share/nginx/html

# Custom nginx config (processed as envsubst template at startup)
COPY frontend/nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]