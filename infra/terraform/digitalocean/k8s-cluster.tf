# Kubernetes cluster configuration for Digital Ocean (optional alternative to droplet)

resource "digitalocean_kubernetes_cluster" "production_k8s" {
  count = var.create_kubernetes_cluster ? 1 : 0
  
  name   = "dansysadm-prod-k8s"
  region = var.region
  version = var.kubernetes_version
  
  # Node pools for production
  node_pool {
    name       = "worker-pool"
    size       = var.k8s_node_size
    count      = var.k8s_node_count
    auto_scale = true
    min_nodes  = var.k8s_min_nodes
    max_nodes  = var.k8s_max_nodes
    
    tags = ["production", "dansysadm.com", "worker"]
  }
  
  # Maintenance window configuration
  maintenance_policy {
    start_time = "03:00"
    day        = "SUNDAY"
  }
  
  # Auto-upgrade settings
  auto_upgrade = true
  
  tags = ["production", "dansysadm.com", "kubernetes"]
}

# Kubernetes node pool for control plane (managed by DO)
resource "digitalocean_kubernetes_node_pool" "control_plane" {
  count = var.create_kubernetes_cluster ? 1 : 0
  
  cluster_id = digitalocean_kubernetes_cluster.production_k8s[0].id
  name       = "control-plane"
  size       = "s-2vcpu-4gb"
  count      = 3
  
  tags = ["production", "dansysadm.com", "control-plane"]
}
