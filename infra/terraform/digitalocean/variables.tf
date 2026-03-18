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
  description = "Digital Ocean SSH key ID for accessing the server (only used if sync_github_admin_keys is false)"
  type        = string
  sensitive   = true
  
  validation {
    condition     = var.sync_github_admin_keys == false ? var.ssh_key_id != "" : true
    error_message = "SSH key ID must be provided when not syncing from GitHub."
  }
}

variable "sync_github_admin_keys" {
  description = "Whether to sync SSH public keys from GitHub user daniellawrence for admin access"
  type        = bool
  default     = true
  
  validation {
    condition     = var.sync_github_admin_keys == true ? var.ssh_key_id == "" : true
    error_message = "SSH key ID should be empty when syncing from GitHub."
  }
}

variable "admin_ssh_ips" {
  description = "List of IP addresses allowed to access SSH (port 22). If empty, will auto-detect from httpbin.org"
  type        = list(string)
  default     = [] # Auto-detected via data source if empty
  
  validation {
    condition     = length(var.admin_ssh_ips) > 0 || var.sync_github_admin_keys == false ? true : true
    error_message = "At least one admin SSH IP must be specified or auto-detected."
  }
}

variable "kubernetes_client_ips" {
  description = "List of IP addresses allowed to access Kubernetes API (port 6443) - restricted to admin IPs for ansible-configured k8s"
  type        = list(string)
  default     = [] # Empty; uses var.admin_ssh_ips instead
  
  validation {
    condition     = true # Always valid as it's optional
    error_message = "Kubernetes client IP configuration is optional."
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
