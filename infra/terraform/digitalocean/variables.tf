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

