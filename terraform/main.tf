terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "2.94.0"
    }
  }
  backend "azurerm" {
    resource_group_name  = "tfstate"
    storage_account_name = "fichistfstate"
    container_name       = "tfstate"
    key                  = "terraform.tfstate"

  }
}

provider "azurerm" {
  subscription_id = "58731601-e920-43f2-af4b-f1a79a8ade9f"
  features {}
}


resource "azurerm_resource_group" "fichis_rg" {
  name     = var.resource_group_name
  location = var.azure_region

}

resource "azurerm_virtual_network" "fichis_vnet" {
  name                = var.vnet_name
  location            = var.azure_region
  resource_group_name = var.resource_group_name
  address_space       = var.vnet_address_ranges

  dynamic "subnet" {
    for_each = toset(var.vnet_subnets)
    iterator = subnet

    content {
      name           = subnet.value.name
      address_prefix = subnet.value.prefix
    }
  }
}
## Storage resources
resource "azurerm_storage_account" "fichis_storage" {
  name                     = var.storage_account_name
  resource_group_name      = var.resource_group_name
  location                 = var.azure_region
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "redis_share" {
  name                 = var.file_share_name
  storage_account_name = var.storage_account_name
  quota                = 2
}


# Container resources
resource "azurerm_container_group" "fichis_container_group" {
  name                = "fichis-cont-group"
  location            = var.azure_region
  resource_group_name = var.resource_group_name

  ip_address_type = "public"
  dns_name_label  = "fichis"
  os_type         = "Linux"
  exposed_port = [
    {
      port     = var.http_port
      protocol = "TCP"
    },
    {
      port     = var.https_port
      protocol = "TCP"
  }]

  container {
    name   = "fichis-api"
    image  = "docker.io/b1t3x/fichis-go"
    cpu    = "0.5"
    memory = "0.5"

    ports {
      port     = var.http_port
      protocol = "TCP"
    }
    ports {
      port     = var.https_port
      protocol = "TCP"
    }

    environment_variables = {
      FICHIS_HTTPS_PORT            = var.https_port
      FICHIS_HTTP_PORT             = var.http_port
      FICHIS_TLS_ON                = var.tls_enabled ? "yes" : "no"
      FICHIS_REDIS_HOST            = var.redis_host
      FICHIS_REDIS_PORT            = var.redis_port
      FICHIS_CERTIFICATE_FILE_PATH = "/app/tls/certificate.crt"
      FICHIS_KEY_FILE_PATH         = "/app/tls/private.key"
      FICHIS_HEALTH_PROBE_PATH     = var.app_gateway_health_probe_path
    }

    volume {
      name       = "certificate"
      mount_path = "/app/tls"
      secret = {
        "certificate.crt" = base64encode(file(var.certificate_file)),
        "private.key"     = base64encode(file(var.key_file))
      }
    }
  }

  container {
    name   = "fichis-redis"
    image  = "docker.io/redis"
    cpu    = "0.5"
    memory = "0.5"

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

resource "azurerm_public_ip" "app_gateway_public_ip" {
  name                = "${var.app_gateway_name}-pip"
  resource_group_name = var.resource_group_name
  location            = var.azure_region
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_application_gateway" "fichis_app_gw" {
  name                = var.app_gateway_name
  resource_group_name = var.resource_group_name
  location            = var.azure_region

  sku {
    name     = "Standard_v2"
    tier     = "Standard_V2"
    capacity = 1
  }
  gateway_ip_configuration {
    name      = "${var.app_gateway_name}-ip-conf"
    subnet_id = (azurerm_virtual_network.fichis_vnet.subnet[*].id)[0]
  }

  frontend_port {
    name = "fich.is-https-port"
    port = 443
  }

  frontend_port {
    name = "fich.is-http-port"
    port = 80
  }

  frontend_ip_configuration {
    name                 = "${var.app_gateway_name}-fe-conf"
    public_ip_address_id = azurerm_public_ip.app_gateway_public_ip.id
  }

  backend_address_pool {
    name  = "fichis-be-pool"
    fqdns = [azurerm_container_group.fichis_container_group.fqdn]
  }

  backend_http_settings {
    name                  = "fichis_http_setting"
    cookie_based_affinity = "Disabled"
    port                  = 443
    protocol              = "Https"
    request_timeout       = 10
    host_name             = "fich.is"
    probe_name            = "fich.is-probe"
  }

  ssl_certificate {
    name     = "fich.is-ssl-certificate"
    data    = filebase64("~/Downloads/fich.is/fichis.pfx")
    password = "Madua123"
  }

  http_listener {
    name                           = "fich.is-https"
    frontend_ip_configuration_name = "${var.app_gateway_name}-fe-conf"
    frontend_port_name             = "fich.is-https-port"
    protocol                       = "Https"

    ssl_certificate_name = "fich.is-ssl-certificate"
  }

  http_listener {
    name                           = "fich.is-http"
    frontend_ip_configuration_name = "${var.app_gateway_name}-fe-conf"
    frontend_port_name             = "fich.is-http-port"
    protocol                       = "Http"
  }

  probe {
    interval                                  = 15
    name                                      = "fich.is-probe"
    protocol                                  = "Https"
    path                                      = var.app_gateway_health_probe_path
    timeout                                   = 45
    unhealthy_threshold                       = 3
    pick_host_name_from_backend_http_settings = true
  }

  request_routing_rule {
    name                       = "fichis-https-routing-rule"
    rule_type                  = "Basic"
    http_listener_name         = "fich.is-https"
    backend_address_pool_name  = "fichis-be-pool"
    backend_http_settings_name = "fichis_http_setting"
  }

  request_routing_rule {
    name                       = "fichis-http-routing-rule"
    rule_type                  = "Basic"
    http_listener_name         = "fich.is-http"
    backend_address_pool_name  = "fichis-be-pool"
    backend_http_settings_name = "fichis_http_setting"
  }

}
## DNS resources

resource "azurerm_dns_zone" "fichis_dns_zone" {
  name                = "fich.is"
  resource_group_name = var.resource_group_name

}

resource "azurerm_dns_a_record" "primary_a_dns_record" {
  name                = "@"
  zone_name           = azurerm_dns_zone.fichis_dns_zone.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [azurerm_public_ip.app_gateway_public_ip.ip_address]
}
