#!/usr/bin/env bash
set -euo pipefail

mkdir -p rootfs
docker pull alpine:latest
ctr=$(docker create alpine:latest /bin/true)
docker export "$ctr" | tar -C rootfs -xvf -
docker rm "$ctr"
echo "✅ rootfs pronta in ./rootfs"

