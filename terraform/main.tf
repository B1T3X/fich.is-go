terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "2.94.0"
    }
  }
  backend "azurerm" {
      resource_group_name = "tfstate"
      storage_account_name = "fichistfstate"
      container_name = "tfstate"
      key = "terraform.tfstate"
    
  }
}

provider "azurerm" {
  features {}
}


resource "azurerm_resource_group" "fichis_rg" {
  name = var.resource_group_name
  location = var.azure_region
}


## Storage resources
resource "azurerm_storage_account" "fichis_storage" {
  name                = var.storage_account_name
  resource_group_name = var.resource_group_name
  location = var.azure_region
  account_tier = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "redis_share" {
    name = var.file_share_name
    storage_account_name = var.storage_account_name
    quota = 2
}


# Container resources
resource "azurerm_container_group" "fichis_container_group" {
  name                = "fichis-cont-group"
  location            = var.azure_region
  resource_group_name = var.resource_group_name

  ip_address_type = "public"
  dns_name_label  = "fichis"
  os_type = "Linux"
  exposed_port    = [
      {
      port = var.http_port
      protocol = "TCP"
  },
  {
      port = var.https_port
      protocol = "TCP"
  }]

  container {
    name   = "fichis-api"
    image  = "docker.io/b1t3x/fichis-go"
    cpu    = "1"
    memory = "1"

    ports {
      port     = var.http_port
      protocol = "TCP"
    }
    ports {
        port = var.https_port
        protocol = "TCP"
    }

    environment_variables = {
      FICHIS_HTTPS_PORT = var.https_port
      FICHIS_HTTP_PORT = var.http_port
      FICHIS_TLS_ON     = var.tls_enabled ? "yes" : "no"
      FICHIS_REDIS_HOST = var.redis_host
      FICHIS_REDIS_PORT = var.redis_port
      FICHIS_CERTIFICATE_FILE_PATH = "/app/tls/certificate.crt"
      FICHIS_KEY_FILE_PATH = "/app/tls/private.key"
    }

    volume {
      name = "certificate"
      mount_path = "/app/tls"
      secret = {
          "certificate.crt" = base64encode(file(var.certificate_file)),
          "private.key" = base64encode(file(var.key_file))
      }
    }
  }

  container {
    name   = "fichis-redis"
    image  = "docker.io/redis"
    cpu    = "1"
    memory = "1"

    ports {
      port     = var.redis_port
      protocol = "TCP"
    }

    volume {
      name                 = "redisdata"
      mount_path           = "/data"
      storage_account_name = var.storage_account_name
      storage_account_key  = azurerm_storage_account.fichis_storage.primary_access_key
      share_name           = var.file_share_name
    }
  }
}

## DNS resources

resource "azurerm_dns_zone" "fichis_dns_zone" {
    name = "fich.is"
    resource_group_name = var.resource_group_name
}

resource "azurerm_dns_a_record" "primary_a_dns_record" {
    name = "@"
    zone_name = azurerm_dns_zone.fichis_dns_zone.name
    resource_group_name = var.resource_group_name
    ttl = 300
    records = [azurerm_container_group.fichis_container_group.ip_address]
}