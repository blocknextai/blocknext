#!/bin/sh

set -e

cat > /usr/share/nginx/html/env.js <<EOF
window.__ENV__ = {
  VITE_PLATFORM_API_URL: "${VITE_PLATFORM_API_URL}",
  VITE_MCP_API_URL: "${VITE_MCP_API_URL}",
}
EOF
