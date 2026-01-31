# TradingAgents Golang Makefile
# ==============================

# 变量定义
APP_NAME := tradingagents
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Go 相关
GO := go
GOFLAGS := -v
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# 目录
BUILD_DIR := ./build
CMD_DIR := ./cmd
INTERNAL_DIR := ./internal

# Docker 相关
DOCKER_IMAGE := tradingagents
DOCKER_TAG := $(VERSION)
DOCKER_REGISTRY ?= 

# Kubernetes 相关
K8S_NAMESPACE := tradingagents
HELM_RELEASE := tradingagents

# 颜色输出
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

.PHONY: all build clean test lint fmt vet help docker-build docker-push k8s-deploy

# 默认目标
all: lint test build

## 帮助信息
help:
	@echo "$(GREEN)TradingAgents Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Usage:$(NC)"
	@echo "  make <target>"
	@echo ""
	@echo "$(YELLOW)Targets:$(NC)"
	@echo "  $(GREEN)build$(NC)          - Build the application"
	@echo "  $(GREEN)build-all$(NC)      - Build for all platforms"
	@echo "  $(GREEN)clean$(NC)          - Clean build artifacts"
	@echo "  $(GREEN)test$(NC)           - Run all tests"
	@echo "  $(GREEN)test-unit$(NC)      - Run unit tests"
	@echo "  $(GREEN)test-integration$(NC) - Run integration tests"
	@echo "  $(GREEN)test-coverage$(NC)  - Run tests with coverage"
	@echo "  $(GREEN)lint$(NC)           - Run linters"
	@echo "  $(GREEN)fmt$(NC)            - Format code"
	@echo "  $(GREEN)vet$(NC)            - Run go vet"
	@echo "  $(GREEN)tidy$(NC)           - Tidy go modules"
	@echo "  $(GREEN)docker-build$(NC)   - Build Docker image"
	@echo "  $(GREEN)docker-push$(NC)    - Push Docker image"
	@echo "  $(GREEN)k8s-deploy$(NC)     - Deploy to Kubernetes"
	@echo "  $(GREEN)k8s-delete$(NC)     - Delete from Kubernetes"
	@echo "  $(GREEN)helm-install$(NC)   - Install via Helm"
	@echo "  $(GREEN)helm-upgrade$(NC)   - Upgrade via Helm"
	@echo "  $(GREEN)helm-uninstall$(NC) - Uninstall via Helm"

## 构建
build:
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)$(NC)"

build-linux:
	@echo "$(GREEN)Building for Linux...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(CMD_DIR)/main.go

build-darwin:
	@echo "$(GREEN)Building for macOS...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(CMD_DIR)/main.go
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(CMD_DIR)/main.go

build-windows:
	@echo "$(GREEN)Building for Windows...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(CMD_DIR)/main.go

build-all: build-linux build-darwin build-windows
	@echo "$(GREEN)All platforms built successfully!$(NC)"

## 清理
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)Clean complete$(NC)"

## 测试
test: test-unit test-integration
	@echo "$(GREEN)All tests passed!$(NC)"

test-unit:
	@echo "$(GREEN)Running unit tests...$(NC)"
	$(GO) test -v -race -short ./internal/...

test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	$(GO) test -v -race -run Integration ./tests/...

test-coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

test-bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	$(GO) test -v -bench=. -benchmem ./...

## 代码质量
lint:
	@echo "$(GREEN)Running linters...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed, running basic checks...$(NC)"; \
		$(GO) vet ./...; \
		$(MAKE) fmt-check; \
	fi

fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GO) fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

fmt-check:
	@echo "$(GREEN)Checking code format...$(NC)"
	@test -z "$$(gofmt -l .)" || (echo "$(RED)Code not formatted. Run 'make fmt'$(NC)" && exit 1)

vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GO) vet ./...

## 依赖管理
tidy:
	@echo "$(GREEN)Tidying go modules...$(NC)"
	$(GO) mod tidy
	$(GO) mod verify

deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GO) mod download

## Docker
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-push:
	@echo "$(GREEN)Pushing Docker image...$(NC)"
ifdef DOCKER_REGISTRY
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest
else
	@echo "$(YELLOW)DOCKER_REGISTRY not set, pushing to default registry$(NC)"
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest
endif

docker-run:
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker run --rm -it \
		-e OPENAI_API_KEY=${OPENAI_API_KEY} \
		-e ALPHA_VANTAGE_API_KEY=${ALPHA_VANTAGE_API_KEY} \
		$(DOCKER_IMAGE):$(DOCKER_TAG) run --ticker NVDA --date 2024-05-10

## Kubernetes
k8s-deploy:
	@echo "$(GREEN)Deploying to Kubernetes...$(NC)"
	kubectl apply -f deploy/k8s/namespace.yaml
	kubectl apply -f deploy/k8s/configmap.yaml
	kubectl apply -f deploy/k8s/secret.yaml
	kubectl apply -f deploy/k8s/deployment.yaml
	kubectl apply -f deploy/k8s/service.yaml
	@echo "$(GREEN)Deployment complete!$(NC)"

k8s-delete:
	@echo "$(YELLOW)Deleting from Kubernetes...$(NC)"
	kubectl delete -f deploy/k8s/ --ignore-not-found
	@echo "$(GREEN)Deletion complete$(NC)"

k8s-status:
	@echo "$(GREEN)Kubernetes status:$(NC)"
	kubectl get all -n $(K8S_NAMESPACE)

k8s-logs:
	kubectl logs -f -l app=$(APP_NAME) -n $(K8S_NAMESPACE)

## Helm
helm-install:
	@echo "$(GREEN)Installing via Helm...$(NC)"
	helm install $(HELM_RELEASE) ./deploy/helm/tradingagents \
		--namespace $(K8S_NAMESPACE) \
		--create-namespace

helm-upgrade:
	@echo "$(GREEN)Upgrading via Helm...$(NC)"
	helm upgrade $(HELM_RELEASE) ./deploy/helm/tradingagents \
		--namespace $(K8S_NAMESPACE)

helm-uninstall:
	@echo "$(YELLOW)Uninstalling via Helm...$(NC)"
	helm uninstall $(HELM_RELEASE) --namespace $(K8S_NAMESPACE)

helm-template:
	@echo "$(GREEN)Rendering Helm templates...$(NC)"
	helm template $(HELM_RELEASE) ./deploy/helm/tradingagents

## 开发工具
install-tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

run:
	@echo "$(GREEN)Running application...$(NC)"
	$(GO) run $(CMD_DIR)/main.go run --ticker NVDA --date 2024-05-10

run-debug:
	@echo "$(GREEN)Running application in debug mode...$(NC)"
	$(GO) run $(CMD_DIR)/main.go run --ticker NVDA --date 2024-05-10 --debug
