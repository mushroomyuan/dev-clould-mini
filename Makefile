PROJECT_NAME := cmdb
OUTPUT_NAME := cmdb
MAIN_FILE := main.go
PKG := github.com/mushroomyuan/dev-clould-mini

gen: ## Init Service
	@protoc -I=. --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG}  cmdb/apps/*/*.proto
	@go fmt ./...
	@protoc-go-inject-tag -input='cmdb/apps/*/*.pb.go'
	@mcube enum -p -m cmdb/apps/*/*.pb.go

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
