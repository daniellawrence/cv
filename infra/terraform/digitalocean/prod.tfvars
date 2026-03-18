# Production Terraform variables for dansysadm.com (Digital Ocean)
# This file contains production-specific configurations

region                      = "sfo1"
droplet_size               = "s-1vcpu-2gb"
image                      = "debian-13-x64"
enable_floating_ip         = true
domain_name                = "dansysadm.com"

# IMPORTANT: Update these with your actual values before running terraform apply
ssh_key_id                 = "" # Leave empty if syncing from GitHub
admin_ssh_ips              = []   # Empty to auto-detect via httpbin.org, or specify IPs manually
kubernetes_client_ips      = [] # Empty if not using K8s

# For production, consider creating a managed Kubernetes cluster:
create_kubernetes_cluster  = false

# Sync SSH keys from GitHub daniellawrence account to DigitalOcean
sync_from_github           = true
