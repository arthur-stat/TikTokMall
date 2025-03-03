#!/bin/bash

log() {
	echo "[INFO] $1"
}

cd ..

services=("auth" "cart" "checkout" "order" "payment" "product" "user")

log "Starting go mod tidy for all services..."

# 遍历所有服务
for service in "${services[@]}"; do
	log "Tidying dependencies for $service..."
	cd app/$service || exit 1
	go mod tidy
	cd - >/dev/null || exit 1
done

# Tidy rpc_gen
log "Tidying dependencies for rpc_gen..."
cd rpc_gen || exit 1
go mod tidy
cd - >/dev/null || exit 1

log "All dependencies tidied!"
