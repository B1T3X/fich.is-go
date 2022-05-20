terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "4.21.0"
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
  features {}
}

provider "google" {
  project = "fichis-go"
  region = var.fichis_google_project_location
}


resource "google_project" "fichis_google_project" {
  name     = var.fichis_google_project_name
  project_id = var.fichis_google_project_id
}

resource "google_service_account" "fichis_google_service_account" {
  account_id   = "fichis-cloudrun-sa"
  display_name = "Cloud Run SA for fich.is"
}

resource "google_cloud_run_service" "fichis_google_cloud_run_service" {
  name     = var.fichis_google_cloud_run_name
  location = var.fichis_google_project_location

  template {
    spec {
      containers {
        image = "europe-west4-docker.pkg.dev/fichis-go/fichis-repo/fichis-go"
        dynamic "env" {
          for_each = var.fichis_google_cloud_run_environment_variables
          iterator = env_var
        content {
          name = env_var.key
          value = env_var.value
}

}
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_domain_mapping" "fichis_google_cloud_run_domain_mapping" {
  location = var.fichis_google_project_location
  name     = var.fichis_custom_domain_name

  metadata {
    namespace = var.fichis_google_project_name
  }

  spec {
    route_name = google_cloud_run_service.fichis_google_cloud_run_service.name
  }
}

resource "google_secret_manager_secret" "fichis_google_secret" {
  secret_id = "fichis_sa_config"

  replication {
    automatic = true
  }
}


resource "google_secret_manager_secret_version" "fichis_google_secret_version" {
  secret = google_secret_manager_secret.fichis_google_secret.id

  secret_data = "secret-data"
}

// TODO: change
data "google_service_account_id_token" "oidc" {
    target_audience = "https://foo.bar/"
  }

  output "oidc_token" {
    value = data.google_service_account_id_token.oidc.id_token
  }
