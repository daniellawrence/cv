HELM_VERSION := v3.14.0
HELMFILE_VERSION := v3.14.0
K3D_VERSION := v5.9.0-rc.0
KUBECTL_VERSION := v1.29.0
KUBECTL_VALIDATE_VERSION := v0.0.4
# Hardcoded for Linux amd64
OS := linux
ARCH := amd64

BUF=buf

.PHONY: setup proto test clean lint provision

install: install-k3d install-kubectl install-helm install-helmfile install-kubectl-validate

provision:
	$(MAKE) -C infra/ansible run

setup:
	go install github.com/bufbuild/buf/cmd/buf@latest

proto: setup
	cd proto && $(BUF) generate

bin/golangci-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.11.2

lint:
	find backend -mindepth 1  -type d -exec ./bin/golangci-lint run {} \;

# k3d
.PHONY: install-k3d
install-k3d: bin/k3d
bin/k3d:
	@echo "Installing k3d $(K3D_VERSION)..."
	@mkdir -p bin/
	@curl -fsSL -o bin/k3d https://github.com/k3d-io/k3d/releases/download/$(K3D_VERSION)/k3d-$(OS)-$(ARCH)
	@chmod +x bin/k3d
	@echo "✓ k3d installed: $$(bin/k3d version)"

# kubectl
.PHONY: install-kubectl
install-kubectl: bin/kubectl
bin/kubectl:
	@echo "Installing kubectl $(KUBECTL_VERSION)..."
	@mkdir -p bin/
	@curl -fsSL -o bin/kubectl https://dl.k8s.io/release/$(KUBECTL_VERSION)/bin/$(OS)/$(ARCH)/kubectl
	@chmod +x bin/kubectl
	@echo "✓ kubectl installed: $$(./bin/kubectl version --client=true|head -1)"

# Helm
.PHONY: install-helm
install-helm: bin/helm
bin/helm:
	@echo "Installing Helm $(HELM_VERSION)..."
	@mkdir -p bin/
	@curl -fsSL -o helm.tar.gz https://get.helm.sh/helm-$(HELM_VERSION)-$(OS)-$(ARCH).tar.gz
	@tar xzf helm.tar.gz -C bin/ --strip-components=1 $(OS)-$(ARCH)/helm
	@rm helm.tar.gz
	@chmod +x bin/helm
	@echo "✓ Helm installed: $$(bin/helm version --short)"

# Helmfile
.PHONY: install-kubectl-validate
install-kubectl-validate: bin/kubectl-validate
bin/kubectl-validate:
	@echo "Installing kubectl-validate $(KUBECTL_VALIDATE_VERSION)..."
	@mkdir -p bin/ ~/.local/bin/
	@curl -sSfL https://github.com/kubernetes-sigs/kubectl-validate/releases/download/$(KUBECTL_VALIDATE_VERSION)/kubectl-validate_$(OS)_$(ARCH).tar.gz \
		| tar -xz -C bin/ kubectl-validate
	@cp bin/kubectl-validate ~/.local/bin/kubectl-validate

.PHONY: install-helmfile
install-helmfile: bin/helmfile
bin/helmfile:
	@echo "Installing Helmfile..."
	@mkdir -p bin/
	@curl -fsSL -o bin/helmfile https://github.com/roboll/helmfile/releases/latest/download/helmfile_$(OS)_$(ARCH)
	@chmod +x bin/helmfile
	@echo "✓ Helmfile installed: $$(bin/helmfile version)"

# Render
.PHONY: render
render:
	./scripts/render.sh