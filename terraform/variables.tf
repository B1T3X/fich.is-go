variable "vnet_name" {
  type    = string
  default = "fichis-app-vnet"
}

variable "vnet_address_ranges" {
  type    = list(string)
  default = ["192.168.0.0/24"]
}

variable "vnet_subnets" {
  type = list(map(string))
  default = [
    {
      name   = "AppSubnet",
      prefix = "192.168.0.0/24"
    }
  ]
}

variable "app_gateway_name" {
  type    = string
  default = "fichis-app-gateway"
}

variable "app_gateway_health_probe_path" {
  type    = string
  default = "/api/healthprobe"
}
variable "container_group_name" {
  type    = string
  default = "fichis-cont-group"
}

variable "storage_account_name" {
  type    = string
  default = "fichisfiles"
}

variable "file_share_name" {
  type    = string
  default = "redisfs"
}

variable "resource_group_name" {
  type    = string
  default = "fichis-app-rg"
}

variable "azure_region" {
  type    = string
  default = "westeurope"
}

variable "tls_enabled" {
  type    = bool
  default = true
}

variable "http_port" {
  type    = number
  default = 80
}

variable "https_port" {
  type    = number
  default = 443
}

variable "redis_host" {
  type    = string
  default = "localhost"
}

variable "redis_port" {
  type    = number
  default = 6379
}

variable "certificate_file" {
  type        = string
  description = "Path to certificate file"
  default     = "~/Downloads/fich.is/certificate.crt"
}

variable "key_file" {
  type        = string
  description = "Path to private key file"
  default     = "~/Downloads/fich.is/private.key"
}
