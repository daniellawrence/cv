# Terraform Infrastructure for Digital Ocean (dansysadm.com)

This directory contains Terraform configuration for managing the production infrastructure on DigitalOcean that hosts **dansysadm.com**.

## Directory Structure

```
infra/terraform/
└── digitalocean/
    ├── main.tf           # Main infrastructure resources
    ├── variables.tf      # Input variable definitions
    ├── output.tf         # Output values
    ├── k8s-cluster.tf    # Optional Kubernetes cluster config
    ├── versions.tf       # Version configurations
    ├── prod.tfvars       # Production-specific variables (update before apply)
    ├── environment.tfvars.example  # Template for local configuration
    └── .gitignore        # Ignore sensitive files
```

## Prerequisites

1. **Terraform** >= 1.0.0
2. **DigitalOcean API Token** - Set as `DO_API_TOKEN` environment variable
3. SSH key created in DigitalOcean (ID needed for configuration)
4. Admin IP addresses for firewall rules

## Setup Instructions

### 1. Configure Environment Variables

```bash
export DO_API_TOKEN="your_digitalocean_api_token"
export GITHUB_TOKEN="your_github_personal_access_token" # For fetching SSH keys from GitHub
```

**Important**: Never commit your API tokens to version control!

### 2. Review and Update Configuration

Copy the example environment file:

```bash
cp environment.tfvars.example terraform.tfvars.local
```

Edit `terraform.tfvars.local` with your specific values:
- SSH key ID from DigitalOcean (or set `sync_from_github = true`)
- Admin IP addresses for firewall access
- Kubernetes configuration (if enabled)

### 3. Initialize Terraform

```bash
cd infra/terraform/digitalocean
terraform init
```

### 4. Plan Changes

```bash
# Review proposed changes before applying
terraform plan -var-file=prod.tfvars
```

### 5. Apply Infrastructure

```bash
terraform apply -var-file=prod.tfvars
```

## SSH Key Synchronization from GitHub

You can automatically sync your GitHub SSH keys to DigitalOcean by setting `sync_from_github = true` in your configuration. This will:

1. Fetch all SSH public keys from your GitHub account (daniellawrence)
2. Create corresponding DigitalOcean SSH key resources
3. Automatically assign the first synced key to the production droplet

**To enable**: Set `sync_from_github = true` in `prod.tfvars` and provide a valid `GITHUB_TOKEN`.

**Note**: When enabled, leave `ssh_key_id` empty as it will be auto-populated from GitHub.

## Automatic IP Detection for SSH Access

The configuration automatically detects your current public IP address using the HTTP data source (https://httpbin.org/ip). This allows you to:

- Set `admin_ssh_ips = []` in your configuration
- The firewall will automatically allow access from your current IP
- Update automatically when you run terraform plan/apply from a different network

**Important**: If you're running Terraform from a different IP than where you'll be administering the server, manually specify your admin IPs instead of relying on auto-detection.

## Resources Managed

- **Droplet**: Production server for dansysadm.com
- **Firewall**: Security rules for SSH, HTTP, HTTPS, and Kubernetes API
- **Floating IP** (optional): Stable public IP address
- **DNS Records**: A record and CNAME for dansysadm.com and www.dansysadm.com

## Optional: Managed Kubernetes Cluster

Set `create_kubernetes_cluster = true` in your variables to deploy a managed DigitalOcean Kubernetes cluster instead of a single droplet. This provides:
- High availability across multiple nodes
- Auto-scaling capabilities
- Automatic upgrades and maintenance

## Security Considerations

1. **Never commit** sensitive data (API tokens, SSH keys) to version control
2. **Restrict SSH access** to known IP addresses only
3. **Enable monitoring** on production resources
4. **Use floating IPs** for stable public endpoints
5. **Regularly rotate** API tokens and SSH keys

## Destroy Infrastructure

To remove all created resources:

```bash
terraform destroy -var-file=prod.tfvars
```

**Warning**: This will delete the production server and all associated resources!

## Troubleshooting

- Check DigitalOcean API token validity
- Verify SSH key ID is correct in your account
- Ensure IP addresses are reachable from your network
- Review Terraform state: `terraform show` or `terraform state list`

## References

- [DigitalOcean Terraform Provider](https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs)
- [Terraform Documentation](https://www.terraform.io/docs/)
- [DigitalOcean Kubernetes Docs](https://docs.digitalocean.com/products/kubernetes/)
