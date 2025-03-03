#!/bin/bash

# 输出日志
log() {
	echo "[INFO] $1"
}

# 输出错误日志
error() {
	echo "[ERROR] $1" >&2
}

# 检查并创建目录
ensure_dir() {
	if [ ! -d "$1" ]; then
		log "Creating directory: $1"
		mkdir -p "$1"
	fi
}

# 生成服务端代码
generate_server_code() {
	local service=$1
	local proto_path=$2

	log "Generating server code for $service..."
	ensure_dir "app/$service"
	cd "app/$service" || exit 1

	cwgo server --type RPC --idl "$proto_path" --server_name "$service" \
		--module "TikTokMall/app/$service" \
		-I ../../idl \
		--idl "../../$proto_path" || error "Failed to generate server code for $service"

	cd - >/dev/null || exit 1
}

# 生成客户端代码
generate_client_code() {
	local service=$1
	local proto_path=$2

	log "Generating client code for $service..."
	ensure_dir "rpc_gen"
	cd "rpc_gen" || exit 1

	cwgo client --type RPC --idl "$proto_path" --server_name "$service" \
		--module "TikTokMall/rpc_gen" \
		-I ../idl \
		--idl "../$proto_path" || error "Failed to generate client code for $service"

	cd - >/dev/null || exit 1
}

# 主函数
main() {
	cd ..

	log "Starting code generation..."

	# 定义服务及其对应的 .proto 文件路径
	declare -A services
	services=(
		["auth"]="idl/auth.proto"
		["cart"]="idl/cart.proto"
		["checkout"]="idl/checkout.proto"
		["order"]="idl/order.proto"
		["payment"]="idl/payment.proto"
		["product"]="idl/product.proto"
		["user"]="idl/user.proto"
	)

	for service in "${!services[@]}"; do
		generate_server_code "$service" "${services[$service]}"
		generate_client_code "$service" "${services[$service]}"
	done

	log "Code generation complete!"
}

# 执行主函数
main
