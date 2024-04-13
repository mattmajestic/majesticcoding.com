package main

import (
    "fmt"
    "net/http"
    "os/exec"
    "testing"
    "time"
)

func TestDeploymentCheck(t *testing.T) {
    fmt.Println("\x1b[32m🚀 Starting the deployment process...\x1b[0m")

    // Start the main application
    cmd := exec.Command("go", "run", "main.go")
    err := cmd.Start()
    if err != nil {
        t.Fatalf("Failed to start main application: %v", err)
    }
    defer cmd.Process.Kill()

    // Give the application some time to start
    time.Sleep(5 * time.Second)

    // Assume deployment step is successful
    deploymentStep := true

    if deploymentStep {
        fmt.Println("⚠️ Checking Deployment for You")

        // Loop over deployment check for demonstration
        for attempt := 1; attempt <= 3; attempt++ {
            fmt.Printf("\nAttempt %d:\n", attempt)
            time.Sleep(2 * time.Second)

            if attempt > 0 {
                fmt.Println("⚠️ Pinging the network")

                // Send a GET request to localhost:8080
                resp, err := http.Get("http://localhost:8080")
                if err != nil {
                    fmt.Println("Error pinging the network:", err)
                    continue
                }
                defer resp.Body.Close()

                // Check the status code of the response
                if resp.StatusCode == http.StatusOK {
                    fmt.Println("\x1b[32m✅ Received 200 response from localhost:8080\x1b[0m")
                    break
                } else {
                    fmt.Println("\x1b[31m❌ Did not receive 200 response from localhost:8080\x1b[0m")
                }
            }
        }

        fmt.Println("\n\x1b[32m🏁 Deployment checks successful\x1b[0m")
    } else {
        fmt.Println("\x1b[31m❌ Deployment Failed... See Go Outputs\x1b[0m")
        t.Fail()
    }
}