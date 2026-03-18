# Terraform configuration for Digital Ocean production host (dansysadm.com)
# Provider: Digital Ocean and HTTP

terraform {
  required_version = ">= 1.0.0"
  
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    
    http = {
      source  = "hashicorp/http"
      version = "~> 4.0"
    }
  }
}

# HTTP provider for fetching dynamic data (public IP, GitHub keys)
provider "http" {}

provider "digitalocean" {
  # Token is loaded from DO_API_TOKEN environment variable
  # Never commit tokens to version control
}

# Fetch your current public IP address from httpbin.org for automatic SSH access
data "http" "my_ip" {
  url = "https://httpbin.org/ip"
  
  request_headers = {
    Accept = "application/json"
  }
}

# Parse the client IP from httpbin response
locals {
  admin_ssh_ips = [jsondecode(data.http.my_ip.response_body).origin]
}

# Use GitHub REST API to fetch SSH keys (more reliable than github provider)
data "http" "github_ssh_keys" {
  url = "https://api.github.com/users/daniellawrence/keys"
  
  request_headers = {
    Accept = "application/vnd.github.v3+json"
  }
}

# Create DigitalOcean SSH key from GitHub's public keys
resource "digitalocean_ssh_key" "from_github" {
  count       = var.sync_from_github ? length(jsondecode(data.http.github_ssh_keys.response_body)) : 0
  
  name        = "daniellawrence-${jsondecode(data.http.github_ssh_keys.response_body)[count.index].title}"
  public_key  = jsondecode(data.http.github_ssh_keys.response_body)[count.index].key
}

# Digital Ocean Droplet (production server for dansysadm.com)
resource "digitalocean_droplet" "production_server" {
  name       = "dansysadm-prod-01"
  region     = var.region
  size       = var.droplet_size
  image      = var.image
  
  # SSH key for access - using GitHub public keys synced to DO or manual ID
  ssh_keys   = var.sync_from_github ? [digitalocean_ssh_key.from_github[0].id] : [var.ssh_key_id]
  
  # Enable monitoring and backups if needed
  monitoring = true
  backups    = false
  
  # Disk resize to 20 GB
  disk       = 20
  
  # Tags for resource management
  tags = ["production", "dansysadm.com", "web-server"]
}

# Firewall to secure the droplet
resource "digitalocean_firewall" "production_firewall" {
  name   = "dansysadm-prod-firewall"
  
  # Allow SSH (for administration) - using auto-detected IP or variable
  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = var.admin_ssh_ips
  }
  
  # Allow HTTP
  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }
  
  # Allow HTTPS (SSL/TLS)
  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }
  
  # Allow Kubernetes API if running K8s cluster
  inbound_rule {
    protocol         = "tcp"
    port_range       = "6443"
    source_addresses = var.kubernetes_client_ips
  }
  
  # Outbound rules (allow all)
  outbound_rule {
    protocol         = "tcp"
    port_range       = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
  
  outbound_rule {
    protocol         = "udp"
    port_range       = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
  
  # Apply to the production droplet
  dropped_traffic_action = "deny"
  droplet_ids            = [digitalocean_droplet.production_server.id]
}

# Floating IP for stable public address (optional)
resource "digitalocean_floating_ip" "production_ip" {
  name       = "dansysadm-prod-fip"
  region     = var.region
  
  # Reserve for production stability
  depends_on = [digitalocean_droplet.production_server]
}

# DNS Record for dansysadm.com
resource "digitalocean_record" "www" {
  domain  = "dansysadm.com"
  type    = "A"
  name    = "@"
  value   = digitalocean_floating_ip.production_ip.ip_address
  ttl     = 300
}

resource "digitalocean_record" "www_www" {
  domain  = "dansysadm.com"
  type    = "CNAME"
  name    = "www"
  value   = "@"
  ttl     = 300
}

# Output your current public IP address (detected from httpbin.org)
output "current_public_ip" {
  description = "Your current public IP address (from httpbin.org)"
  value       = jsondecode(data.http.my_ip.response_body).origin
}

# Output GitHub SSH keys information
output "github_ssh_keys_info" {
  description = "SSH keys from GitHub account daniellawrence"
  value       = jsondecode(data.http.github_ssh_keys.response_body)
}

# Output the synced DigitalOcean SSH key ID for use in droplet configuration
output "synced_digitalocean_key_id" {
  description = "DigitalOcean SSH key ID synced from GitHub (when sync_from_github is true)"
  value       = var.sync_from_github ? digitalocean_ssh_key.from_github[0].id : null
  
  depends_on = [digitalocean_ssh_key.from_github]
}
