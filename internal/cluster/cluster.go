/*
Copyright © 2024 Mathieu DE SOUSA <m.desousa@bl-solutions.co>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cluster

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

func ensureRegistryRunning(verbose bool) error {
    // Check if registry already exists (running or stopped)
    checkExistsCmd := exec.Command("docker", "ps", "-a", "-q", "-f", "name=local-registry")
    existsOutput, _ := checkExistsCmd.Output()
    
    if len(strings.TrimSpace(string(existsOutput))) > 0 {
        // Registry exists, check if it's running
        checkRunningCmd := exec.Command("docker", "ps", "-q", "-f", "name=local-registry")
        runningOutput, _ := checkRunningCmd.Output()
        
        if len(strings.TrimSpace(string(runningOutput))) == 0 {
            // Registry exists but is stopped, start it
            startCmd := exec.Command("docker", "start", "local-registry")
            if verbose {
                startCmd.Stdout = os.Stdout
                startCmd.Stderr = os.Stderr
            }
            if err := startCmd.Run(); err != nil {
                return fmt.Errorf("failed to start existing registry: %w", err)
            }
        }
    } else {
        // Registry doesn't exist, create it
        registryCmd := exec.Command("docker", "run", "-d", "--name", "local-registry", "-p", "5000:5000", "--restart=always", "registry:2")
        
        if verbose {
            registryCmd.Stdout = os.Stdout
            registryCmd.Stderr = os.Stderr
        }
        
        if err := registryCmd.Run(); err != nil {
            return fmt.Errorf("docker registry creation failed: %w", err)
        }
    }
    
    return nil
}

func stopRegistry(verbose bool) error {
    stopCmd := exec.Command("docker", "stop", "local-registry")
    if verbose {
        stopCmd.Stdout = os.Stdout
        stopCmd.Stderr = os.Stderr
    }
    stopCmd.Run() // Ignore errors if registry doesn't exist
    return nil
}

func removeRegistry(verbose bool) error {
    removeCmd := exec.Command("docker", "rm", "local-registry")
    if verbose {
        removeCmd.Stdout = os.Stdout
        removeCmd.Stderr = os.Stderr
    }
    removeCmd.Run() // Ignore errors if registry doesn't exist
    return nil
}

func Create(verbose bool) error {
    // Get user cache directory
    cacheDir, err := os.UserCacheDir()
    if err != nil {
        return fmt.Errorf("failed to get user cache directory: %w", err)
    }

    // Create CLI cache directory if it doesn't exist
    cliCacheDir := filepath.Join(cacheDir, "cli")
    if err := os.MkdirAll(cliCacheDir, 0755); err != nil {
        return fmt.Errorf("failed to create CLI cache directory: %w", err)
    }

    // Generate registry configuration file
    registryConfigPath := filepath.Join(cliCacheDir, "registry.yaml")
    registryConfig := `mirrors:
  "localhost:5000":
    endpoint:
      - "http://localhost:5000"
configs:
  "localhost:5000":
    auth:
      username: ""
      password: ""
    tls:
      insecure_skip_verify: true
`

    if err := os.WriteFile(registryConfigPath, []byte(registryConfig), 0644); err != nil {
        return fmt.Errorf("failed to write registry configuration: %w", err)
    }

    // Ensure Docker registry is running
    if err := ensureRegistryRunning(verbose); err != nil {
        return err
    }

    // Create k3d cluster with registry configuration
    cmd := exec.Command("k3d", "cluster", "create", "local", "--registry-config", registryConfigPath)

    if verbose {
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("k3d cluster creation failed: %w", err)
    }

    return nil
}

func Delete(verbose bool, removeRegistry bool) error {
    // Delete k3d cluster named "local"
    cmd := exec.Command("k3d", "cluster", "delete", "local")

    if verbose {
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
    }

    if err := cmd.Run(); err != nil {
        return fmt.Errorf("k3d cluster deletion failed: %w", err)
    }

    // Stop local Docker registry
    stopRegistry(verbose)

    // Remove registry only if --remove-registry flag is set
    if removeRegistry {
        removeRegistry(verbose)
    }

    return nil
}
