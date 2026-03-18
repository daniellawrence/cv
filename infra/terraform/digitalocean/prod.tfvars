# Production Terraform variables for dansysadm.com (Digital Ocean)
# This file contains production-specific configurations

region                      = "sfo1"
droplet_size               = "s-1vcpu-2gb"
image                      = "debian-13-x64"
enable_floating_ip         = true
domain_name                = "dansysadm.com"

# IMPORTANT: Update these with your actual values before running terraform apply
ssh_key_id                 = "" # Leave empty when syncing from GitHub (recommended)
admin_ssh_ips              = []   # Empty to auto-detect via httpbin.org, or specify IPs manually
kubernetes_client_ips      = []   # Empty; uses var.admin_ssh_ips for k8s API access

# Sync SSH public keys from GitHub user daniellawrence for admin access
sync_github_admin_keys     = true
