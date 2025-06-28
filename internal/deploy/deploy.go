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
)

type AppConfig struct {
	ProjectPath string       `mapstructure:"project_path"`
	Deploy      DeployConfig `mapstructure:"deploy"`
}

type DeployConfig struct {
	ChartPath  string `mapstructure:"chart_path"`
	ValuesFile string `mapstructure:"values_file"`
	Namespace  string `mapstructure:"namespace"`
}

type DependencyConfig struct {
	ChartName   string `mapstructure:"chart_name"`
	ValuesFile  string `mapstructure:"values_file"`
	Version     string `mapstructure:"version"`
	Namespace   string `mapstructure:"namespace"`
	Optional    bool   `mapstructure:"optional,omitempty"`
}

type DependenciesConfig struct {
	Dependencies map[string]DependencyConfig `mapstructure:"dependencies"`
}

func DeployDependencies(config DependenciesConfig, optional bool, verbose bool) error {
	for depName, depConfig := range config.Dependencies {
		// Skip optional dependencies if not requested
		if depConfig.Optional && !optional {
			continue
		}

		// Build Helm command
		args := []string{"upgrade", "--install", depName, depConfig.ChartName}
		
		// Add version if specified
		if depConfig.Version != "" {
			args = append(args, "--version", depConfig.Version)
		}
		
		// Add namespace if specified
		if depConfig.Namespace != "" {
			args = append(args, "--namespace", depConfig.Namespace, "--create-namespace")
		}
		
		// Add values file if specified
		if depConfig.ValuesFile != "" {
			args = append(args, "-f", depConfig.ValuesFile)
		}

		// Execute Helm command
		cmd := exec.Command("helm", args...)
		
		if verbose {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("helm deployment failed for dependency '%s': %w", depName, err)
		}
	}

	return nil
}

func DeployApp(config AppConfig, appName string, verbose bool) error {
	// Validate required fields
	if config.Deploy.ChartPath == "" {
		return fmt.Errorf("chart_path is required")
	}
	if config.Deploy.ValuesFile == "" {
		return fmt.Errorf("values_file is required")
	}
	if config.Deploy.Namespace == "" {
		return fmt.Errorf("namespace is required")
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
	args := []string{"upgrade", "--install", appName, chartPath, "-f", valuesPath, "--namespace", config.Deploy.Namespace, "--create-namespace"}

	// Execute Helm command
	cmd := exec.Command("helm", args...)
	
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("helm deployment failed: %w", err)
	}

	return nil
}