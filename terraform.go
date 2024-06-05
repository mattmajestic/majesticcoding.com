package main

import (
    "fmt"
    "os"
    "os/exec"
)

func main() {
    // Check if Minikube is installed
    if err := runCommand("minikube", "version"); err != nil {
        fmt.Printf("Minikube is not installed or not in PATH: %s\n", err)
        return
    }

    // Start Minikube
    if err := runCommand("minikube", "start"); err != nil {
        fmt.Printf("Failed to start Minikube: %s\n", err)
        return
    }

    // Set Kubernetes context to Minikube
    if err := runCommand("kubectl", "config", "use-context", "minikube"); err != nil {
        fmt.Printf("Failed to set context to Minikube: %s\n", err)
        return
    }

    // Run Terraform init
    if err := runCommand("terraform", "init"); err != nil {
        fmt.Printf("Failed to run 'terraform init': %s\n", err)
        return
    }

    // Run Terraform plan
    if err := runCommand("terraform", "plan"); err != nil {
        fmt.Printf("Failed to run 'terraform plan': %s\n", err)
        return
    }

    // Run Terraform apply
    if err := runCommand("terraform", "apply", "-auto-approve"); err != nil {
        fmt.Printf("Failed to run 'terraform apply': %s\n", err)
        return
    }

    fmt.Println("Kubernetes deployment and service have been successfully created!")
}

func runCommand(command string, args ...string) error {
    cmd := exec.Command(command, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
