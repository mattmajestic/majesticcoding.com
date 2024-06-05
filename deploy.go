package main

import (
    "context"
    "fmt"
    "os"
    "os/exec"

    "google.golang.org/api/container/v1"
    "google.golang.org/api/option"
)

func main() {
    project := os.Getenv("GCP_PROJECT_ID")
    cluster := os.Getenv("GKE_CLUSTER_NAME")
    zone := os.Getenv("GKE_ZONE")
    manifest := os.Getenv("K8S_MANIFEST")

    ctx := context.Background()

    // Initialize GKE client
    containerService, err := container.NewService(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
    if err != nil {
        fmt.Printf("Failed to create GKE client: %s\n", err)
        return
    }

    // Get GKE cluster
    clusterResp, err := containerService.Projects.Zones.Clusters.Get(project, zone, cluster).Context(ctx).Do()
    if err != nil {
        fmt.Printf("Failed to get GKE cluster: %s\n", err)
        return
    }

    // Print cluster information
    fmt.Printf("Cluster Info: %+v\n", clusterResp)

    // Authenticate kubectl to the GKE cluster
    if err := runCommand("gcloud", "auth", "activate-service-account", "--key-file", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")); err != nil {
        fmt.Printf("Failed to authenticate with GCP: %s\n", err)
        return
    }

    // Configure kubectl to use the GKE cluster
    if err := runCommand("gcloud", "container", "clusters", "get-credentials", cluster, "--zone", zone, "--project", project); err != nil {
        fmt.Printf("Failed to configure kubectl: %s\n", err)
        return
    }

    // Apply Kubernetes manifests
    if err := runCommand("kubectl", "apply", "-f", manifest); err != nil {
        fmt.Printf("Failed to apply Kubernetes manifests: %s\n", err)
        return
    }

    fmt.Println("Kubernetes manifests have been successfully applied to GKE!")
}

func runCommand(command string, args ...string) error {
    cmd := exec.Command(command, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
