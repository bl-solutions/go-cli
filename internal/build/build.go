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
package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type BuildConfig struct {
	ProjectPath string            `mapstructure:"project_path"`
	Build       BuildDetails      `mapstructure:"build"`
}

type BuildDetails struct {
	ImageName  string   `mapstructure:"image_name"`
	Dockerfile string   `mapstructure:"dockerfile"`
	Context    string   `mapstructure:"context"`
	BuildArgs  []string `mapstructure:"build_args,omitempty"`
}

func Build(config BuildConfig) error {
	// Validate required fields
	if config.Build.ImageName == "" {
		return fmt.Errorf("image_name is required")
	}
	if config.Build.Dockerfile == "" {
		return fmt.Errorf("dockerfile is required")
	}
	if config.Build.Context == "" {
		return fmt.Errorf("context is required")
	}

	// Change to project directory
	if err := os.Chdir(config.ProjectPath); err != nil {
		return fmt.Errorf("failed to change to project directory '%s': %w", config.ProjectPath, err)
	}

	// Build Docker command
	args := []string{"build"}
	
	// Add tag
	args = append(args, "-t", config.Build.ImageName)
	
	// Add dockerfile path
	dockerfilePath := filepath.Join(config.Build.Context, config.Build.Dockerfile)
	args = append(args, "-f", dockerfilePath)
	
	// Add build args
	for _, buildArg := range config.Build.BuildArgs {
		args = append(args, "--build-arg", buildArg)
	}
	
	// Add context
	args = append(args, config.Build.Context)

	// Execute Docker command
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	return nil
}