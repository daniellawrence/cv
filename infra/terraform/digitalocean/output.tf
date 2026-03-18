# Output values for Digital Ocean production infrastructure (dansysadm.com)

output "droplet_id" {
  description = "ID of the production droplet"
  value       = digitalocean_droplet.production_server.id
}

output "droplet_name" {
  description = "Name of the production droplet"
  value       = digitalocean_droplet.production_server.name
}

output "droplet_ip" {
  description = "Public IP address of the production droplet (use floating IP if enabled)"
  value       = enable_floating_ip ? digitalocean_floating_ip.production_ip.ip_address : digitalocean_droplet.production_server.ipv4_address
}

output "floating_ip" {
  description = "Floating IP address for stable public access"
  value       = enable_floating_ip ? digitalocean_floating_ip.production_ip.ip_address : null
}

output "firewall_id" {
  description = "ID of the production firewall"
  value       = digitalocean_firewall.production_firewall.id
}

output "dns_records" {
  description = "DNS records for dansysadm.com"
  value = {
    www = digitalocean_record.www.fqdn
    a   = digitalocean_record.www.value
  }
}
