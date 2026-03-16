#!/bin/sh
# Replace environment variables in nginx config before starting

set -e

if [ -n "${CONTENT_SHA1}" ] || [ -n "${BUILD_TIMESTAMP}" ]; then
    # Create a temporary nginx config with substituted values
    CONTENT_SHA1_VAL="${CONTENT_SHA1:-none}"
    BUILD_TIMESTAMP_VAL="${BUILD_TIMESTAMP:-unknown}"
    
    envsubst '${CONTENT_SHA1} ${BUILD_TIMESTAMP}' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf
fi
