# TFLint configuration for Terraform validation at root level

config {
  # Fail on error severity only (not warning)
  force = false
  
  # Disable plugin_dir to avoid plugin loading issues
}

# Rule configurations - using built-in rules only (no external plugins)
rule "terraform_deprecated_interpolation" {
  enabled = true
}

rule "terraform_deprecated_index" {
  enabled = true
}

rule "terraform_unused_declarations" {
  enabled = true
}

rule "terraform_comment_syntax" {
  enabled = true
}

rule "terraform_documented_outputs" {
  enabled = false # Optional: require documented outputs
}

rule "terraform_documented_variables" {
  enabled = false # Optional: require documented variables
}

rule "terraform_module_pinned_source" {
  enabled = true
}

rule "terraform_naming_convention" {
  enabled = true
}

# Custom rules for security
rule "terraform_required_providers" {
  enabled = true
}
