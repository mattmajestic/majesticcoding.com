provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "main" {
  name     = "majesticcoding-rg"
  location = "East US"
}

resource "azurerm_container_app_environment" "main" {
  name                = "majesticcoding-env"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
}

resource "azurerm_container_app" "main" {
  name                         = "majesticcoding"
  container_app_environment_id = azurerm_container_app_environment.main.id
  resource_group_name          = azurerm_resource_group.main.name
  revision_mode                = "Single"

  template {
    container {
      name   = "majesticcoding"
      image  = "docker.io/mattmajestic/majesticcoding:latest"
      cpu    = 0.5
      memory = "1.0Gi"
    }
  }

  ingress {
    external_enabled = true
    target_port      = 8080
  }
}