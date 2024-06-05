provider "kubernetes" {
  config_path = "~/.kube/config"
}

resource "null_resource" "apply_manifest" {
  provisioner "local-exec" {
    command = "kubectl apply -f ./k8s-go.yaml"
  }
}
