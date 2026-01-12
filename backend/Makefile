.PHONY: help install run build clean test docker-build docker-run

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## 安装依赖
	go mod download
	go mod tidy

run: ## 运行服务
	go run main.go

build: ## 编译项目
	go build -o bin/auction-backend main.go

clean: ## 清理构建文件
	rm -rf bin/
	rm -f *.log

test: ## 运行测试
	go test -v ./...

db-init: ## 初始化数据库
	mysql -u root -p < schema.sql

docker-build: ## 构建 Docker 镜像
	docker build -t auction-backend:latest .

docker-run: ## 运行 Docker 容器
	docker run -p 8080:8080 --env-file .env auction-backend:latest

dev: ## 开发模式（使用 air 热重载）
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air
