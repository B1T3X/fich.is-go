variable "fichis_google_project_name" {
  type    = string
  default = "fichis-go"
}

variable "fichis_google_project_id" {
  type    = string
  default = "fichis-go"
}

variable "fichis_google_cloud_run_name" {
  type = string
  default = "fichis-go-firestore"
}

variable "fichis_google_project_location" {
  type    = string
  default = "europe-west4"
}

variable "fichis_google_cloud_run_environment_variables" {
  type    = map(string)
  default = {
    "testKey" = "testValue"
    "FICHIS_API_VALIDATION_ON" = "yes"
  }
}
variable "fichis_custom_domain_name" {
  type    = string
  default = "fich.is"
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
