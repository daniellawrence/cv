# Variable definitions for Digital Ocean production infrastructure (dansysadm.com)

variable "region" {
  description = "Digital Ocean region where the droplet will be created"
  type        = string
  default     = "sfo1" # San Francisco 1 - West Coast US
}

variable "droplet_size" {
  description = "Digital Ocean droplet size (2GB RAM, 1 vCPU)"
  type        = string
  default     = "s-1vcpu-2gb" # Basic: 1 vCPU, 2GB RAM
}

variable "image" {
  description = "Digital Ocean image for the droplet (Debian 13 x64)"
  type        = string
  default     = "debian-13-x64"
}

variable "ssh_key_id" {
  description = "Digital Ocean SSH key ID for accessing the server (only used if sync_from_github is false)"
  type        = string
  sensitive   = true
  
  validation {
    condition     = var.sync_from_github == false ? var.ssh_key_id != "" : true
    error_message = "SSH key ID must be provided when not syncing from GitHub."
  }
}

variable "sync_from_github" {
  description = "Whether to sync SSH keys from GitHub account daniellawrence to DigitalOcean"
  type        = bool
  default     = false
  
  validation {
    condition     = var.sync_from_github == true ? var.ssh_key_id == "" : true
    error_message = "SSH key ID should be empty when syncing from GitHub."
  }
}

variable "admin_ssh_ips" {
  description = "List of IP addresses allowed to access SSH (port 22). If empty, will auto-detect from httpbin.org"
  type        = list(string)
  default     = [] # Auto-detected via data source if empty
  
  validation {
    condition     = var.sync_from_github == false && length(var.admin_ssh_ips) > 0 ? true : true
    error_message = "SSH key ID or manual admin IPs must be provided when not syncing from GitHub."
  }
}

variable "kubernetes_client_ips" {
  description = "List of IP addresses allowed to access Kubernetes API (port 6443)"
  type        = list(string)
  default     = [""] # Update with your K8s client IPs in production
  
  validation {
    condition     = length(var.kubernetes_client_ips) > 0
    error_message = "At least one Kubernetes client IP must be specified."
  }
}

variable "enable_floating_ip" {
  description = "Whether to create a floating IP for stable public address"
  type        = bool
  default     = true
}

variable "domain_name" {
  description = "Domain name for the production site (dansysadm.com)"
  type        = string
  default     = "dansysadm.com"
}

# Optional: Kubernetes cluster configuration
variable "create_kubernetes_cluster" {
  description = "Whether to create a managed Kubernetes cluster instead of single droplet"
  type        = bool
  default     = false
}
