#!/bin/sh
set -e

if [ -z "$TARGET_URL" ]; then
  echo "Error: TARGET_URL not set"
  exit 1
fi

echo "Fetching HTML from $TARGET_URL..."
curl -s "$TARGET_URL" -o /usr/share/nginx/html/index.html || echo "<h1>Failed to fetch</h1>" > /usr/share/nginx/html/index.html

echo "Starting nginx..."
exec nginx -g "daemon off;"
