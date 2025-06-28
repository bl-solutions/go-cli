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
package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type AppConfig struct {
	ProjectPath string       `mapstructure:"project_path"`
	Deploy      DeployConfig `mapstructure:"deploy"`
}

type DeployConfig struct {
	ChartPath  string `mapstructure:"chart_path"`
	ValuesFile string `mapstructure:"values_file"`
}

func DeployDependencies(optional bool) error {
	time.Sleep(2 * time.Second)
	return nil
}

func DeployApp(config AppConfig, appName string) error {
	// Validate required fields
	if config.Deploy.ChartPath == "" {
		return fmt.Errorf("chart_path is required")
	}
	if config.Deploy.ValuesFile == "" {
		return fmt.Errorf("values_file is required")
	}

	// Change to project directory
	if err := os.Chdir(config.ProjectPath); err != nil {
		return fmt.Errorf("failed to change to project directory '%s': %w", config.ProjectPath, err)
	}

	// Resolve chart path (relative to project or absolute)
	chartPath := config.Deploy.ChartPath
	if !filepath.IsAbs(chartPath) {
		chartPath = filepath.Join(config.ProjectPath, chartPath)
	}

	// Resolve values file path (relative to project or absolute)
	valuesPath := config.Deploy.ValuesFile
	if !filepath.IsAbs(valuesPath) {
		valuesPath = filepath.Join(config.ProjectPath, valuesPath)
	}

	// Build Helm command
	args := []string{"upgrade", "--install", appName, chartPath, "-f", valuesPath}

	// Execute Helm command
	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("helm deployment failed: %w", err)
	}

	return nil
}