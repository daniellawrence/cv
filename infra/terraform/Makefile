.PHONY: lint init plan apply

# Use binaries from git-root/bin directory
TFLINT := ../../bin/tflint
TERRAFORM := ../../bin/terraform

# Lint Terraform files with tflint
lint:
	$(TFLINT) --init && $(TFLINT)

# Initialize Terraform with backend configuration and install plugins
init:
	$(TERRAFORM) init -reconfigure

# Generate execution plan for Terraform changes
plan:
	$(TERRAFORM) plan

# Apply Terraform changes to infrastructure
apply:
	$(TERRAFORM) apply -auto-approve
