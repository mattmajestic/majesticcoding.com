package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	applyInfra()
}

func applyInfra() {
	steps := []struct {
		name     string
		command  string
		args     []string
		useTFDir bool
	}{
		{"Checking Minikube", "minikube", []string{"version"}, false},
		{"Checking Terraform", "terraform", []string{"version"}, true},
		{"Starting Minikube", "minikube", []string{"start"}, false},
		{"Set Kube Context", "kubectl", []string{"config", "use-context", "minikube"}, false},
		{"Terraform Init", "terraform", []string{"init"}, true},
		{"Terraform Plan", "terraform", []string{"plan"}, true},
		{"Terraform Apply", "terraform", []string{"apply", "-auto-approve"}, true},
	}

	for _, step := range steps {
		fmt.Println("==>", step.name)
		if err := runCommand(step.command, step.args, step.useTFDir); err != nil {
			fmt.Printf("❌ %s failed: %s\n", step.name, err)
			return
		}
	}
	fmt.Println("✅ Infrastructure applied successfully!")
}

func runCommand(command string, args []string, useTFDir bool) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if useTFDir {
		cmd.Dir = "./terraform"
	}
	return cmd.Run()
}
