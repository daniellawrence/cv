# Version variables for Kubernetes configuration

variable "kubernetes_version" {
  description = "Kubernetes version for the managed cluster"
  type        = string
  default     = "1.28" # Latest stable LTS version
}

variable "k8s_node_size" {
  description = "Size of worker nodes in Kubernetes cluster"
  type        = string
  default     = "g-4vcpu-8gb"
}

variable "k8s_node_count" {
  description = "Initial number of worker nodes"
  type        = number
  default     = 3
}

variable "k8s_min_nodes" {
  description = "Minimum number of nodes during auto-scaling"
  type        = number
  default     = 2
}

variable "k8s_max_nodes" {
  description = "Maximum number of nodes during auto-scaling"
  type        = number
  default     = 10
}
