provider "kubernetes" {
  config_path = "~/.kube/config"
}

resource "kubernetes_deployment" "majesticcoding" {
  metadata {
    name = "majesticcoding-deployment"
    labels = {
      app = "majesticcoding"
    }
  }

  spec {
    replicas = 1

    selector {
      match_labels = {
        app = "majesticcoding"
      }
    }

    template {
      metadata {
        labels = {
          app = "majesticcoding"
        }
      }

      spec {
        container {
          name  = "majesticcoding"
          image = "mattmajestic/majesticcoding:latest"

          port {
            container_port = 8080
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "majesticcoding" {
  metadata {
    name = "majesticcoding-service"
  }

  spec {
    selector = {
      app = "majesticcoding"
    }

    port {
      protocol    = "TCP"
      port        = 80
      target_port = 8080
    }

    type = "LoadBalancer"
  }
}
