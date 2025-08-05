package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func waitForServerReady(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("server not ready at %s", url)
}

func TestMainDeployment(t *testing.T) {
	fmt.Println("🚀 Starting the deployment test...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil                                       // important
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // detach

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start app: %v", err)
	}

	defer func() {
		fmt.Println("🛑 Cleaning up app process...")
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	if err := waitForServerReady("http://localhost:8080", 10*time.Second); err != nil {
		t.Fatalf("❌ Failed to connect: %v", err)
	}

	resp, err := http.Get("http://localhost:8080/swagger/index.html")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("❌ Swagger check failed: %v", err)
	}
	resp.Body.Close()

	fmt.Println("✅ Server responded correctly, test complete.")
}
