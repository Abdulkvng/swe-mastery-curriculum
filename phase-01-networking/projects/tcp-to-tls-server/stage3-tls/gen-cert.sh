#!/usr/bin/env bash
# gen-cert.sh - Generate a self-signed TLS cert for testing.
#
# Why self-signed? We don't have a real CA-signed cert for "localhost".
# A self-signed cert means we sign our own cert with our own key.
# Browsers/curl will warn — that's why curl needs -k.
#
# In production at Datadog/Apple, certs are issued by an internal CA
# (or Let's Encrypt for public-facing services).

set -euo pipefail

cd "$(dirname "$0")"

if [[ -f cert.pem && -f key.pem ]]; then
    echo "cert.pem and key.pem already exist. Delete them to regenerate."
    exit 0
fi

# Generate a 2048-bit RSA key + a self-signed cert valid for 365 days.
# CN (Common Name) = "localhost" so it matches when we connect to localhost.
# subjectAltName = also accept "localhost" via SAN, which modern clients require.
openssl req -x509 -newkey rsa:2048 -nodes \
    -keyout key.pem \
    -out cert.pem \
    -days 365 \
    -subj "/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

echo "Generated:"
ls -lh cert.pem key.pem

echo
echo "Inspect the cert:"
echo "  openssl x509 -in cert.pem -text -noout"
