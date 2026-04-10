variable "region" {
  description = "AWS Region."
  type        = string
}

variable "eks_component_name" {
  type        = string
  description = "The name of the eks component"
  default     = "eks/cluster"
}

variable "is_default" {
  type        = bool
  description = "Set `true` to make this the default IngressClass. There should only be one default per cluster."
  default     = false
}

variable "class_name" {
  type        = string
  description = "Class name for default ingress"
  default     = "default"
}

variable "group" {
  type        = string
  description = "Group name for default ingress"
  default     = "common"
}

variable "scheme" {
  type        = string
  description = "Scheme for default ingress, one of `internet-facing` or `internal`."
  default     = "internet-facing"

  validation {
    condition     = contains(["internet-facing", "internal"], var.scheme)
    error_message = "The default ingress scheme must be one of `internet-facing` or `internal`."
  }
}

variable "ip_address_type" {
  type        = string
  description = "IP address type for default ingress, one of `ipv4` or `dualstack`."
  default     = "dualstack"

  validation {
    condition     = contains(["ipv4", "dualstack"], var.ip_address_type)
    error_message = "The default ingress IP address type must be one of `ipv4` or `dualstack`."
  }
}

variable "load_balancer_attributes" {
  type        = list(object({ key = string, value = string }))
  description = <<-EOT
    A list of load balancer attributes to apply to the default ingress load balancer.
    See [Load Balancer Attributes](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/application-load-balancers.html#load-balancer-attributes).
    EOT
  default     = []
}

variable "additional_tags" {
  type        = map(string)
  description = "Additional tags to apply to the ingress load balancer."
  default     = {}
}

variable "auto_mode_enabled" {
  type        = bool
  description = <<-EOT
    Enable EKS Auto Mode. When enabled, AWS manages compute (via managed Karpenter),
    networking (elastic load balancing), and storage (EBS block storage) for the cluster.
    Requires Kubernetes 1.29+ and AWS provider >= 5.79.0.
    Cannot be used with `karpenter_iam_role_enabled = true` (Auto Mode includes managed Karpenter).
    When enabled, the addons `vpc-cni`, `kube-proxy`, `coredns`, and `aws-ebs-csi-driver` are
    managed by Auto Mode and should be removed from the `addons` variable (or use `auto_mode_upgrade`
    for brownfield migration).
    EOT
  default     = false
  nullable    = false
}
