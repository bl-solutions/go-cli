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
package helm

import (
    "fmt"
    "os"
    "os/exec"
)

type RepoConfig struct {
    URL string `mapstructure:"url"`
}

func ConfigureRepos(repos map[string]RepoConfig, verbose bool) error {
    for repoName, repoConfig := range repos {
        // Build Helm repo add command
        args := []string{"repo", "add", repoName, repoConfig.URL}
        
        // Execute Helm command
        cmd := exec.Command("helm", args...)
        
        if verbose {
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
        }
        
        if err := cmd.Run(); err != nil {
            return fmt.Errorf("helm repo add failed for repository '%s': %w", repoName, err)
        }
    }
    
    // Update repositories
    updateCmd := exec.Command("helm", "repo", "update")
    
    if verbose {
        updateCmd.Stdout = os.Stdout
        updateCmd.Stderr = os.Stderr
    }
    
    if err := updateCmd.Run(); err != nil {
        return fmt.Errorf("helm repo update failed: %w", err)
    }
    
    return nil
}