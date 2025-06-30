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
package uninstall

import (
    "fmt"
    "time"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/briandowns/spinner"
    "go-cli/internal/deploy"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
    Use:   "uninstall",
    Short: "Uninstall resources from the cluster",
    Long:  `Uninstall various resources like dependencies from the cluster.`,
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
    },
}

// dependencyCmd represents the dependency subcommand
var dependencyCmd = &cobra.Command{
    Use:   "dependency [dependency-name]",
    Short: "Uninstall a specific dependency",
    Long:  `Uninstall a specific dependency like PostgreSQL, Redis, etc. from the cluster.`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        depName := args[0]
        verbose, _ := cmd.Flags().GetBool("verbose")
        
        // Read dependencies configuration
        var deps map[string]deploy.DependencyConfig
        if err := viper.UnmarshalKey("dependencies", &deps); err != nil {
            fmt.Printf("Error reading dependencies configuration: %v\n", err)
            return
        }
        
        // Check if dependency exists in configuration
        depConfig, exists := deps[depName]
        if !exists {
            fmt.Printf("Dependency '%s' not found in configuration\n", depName)
            return
        }
        
        var s *spinner.Spinner
        if !verbose {
            s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
            s.Suffix = fmt.Sprintf(" Uninstalling dependency %s...", depName)
            s.Start()
        }
        
        err := deploy.UninstallDependency(depName, depConfig, verbose)
        
        if !verbose && s != nil {
            s.Stop()
        }
        
        if err != nil {
            fmt.Printf("Error uninstalling dependency '%s': %v\n", depName, err)
        } else {
            fmt.Printf("Dependency '%s' uninstalled successfully!\n", depName)
        }
    },
}

// appCmd represents the app subcommand for uninstall
var appCmd = &cobra.Command{
    Use:   "app [app-name]",
    Short: "Uninstall an application",
    Long:  `Uninstall a specific application from the cluster.`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        appName := args[0]
        verbose, _ := cmd.Flags().GetBool("verbose")
        
        // Read configuration for the application
        var config deploy.AppConfig
        configKey := fmt.Sprintf("apps.%s", appName)
        if err := viper.UnmarshalKey(configKey, &config); err != nil {
            fmt.Printf("Error reading configuration for app '%s': %v\n", appName, err)
            return
        }
        
        // Check if configuration exists
        if config.ProjectPath == "" {
            fmt.Printf("No configuration found for application '%s'\n", appName)
            return
        }
        
        var s *spinner.Spinner
        if !verbose {
            s = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
            s.Suffix = fmt.Sprintf(" Uninstalling application %s...", appName)
            s.Start()
        }
        
        err := deploy.UninstallApp(config, appName, verbose)
        
        if !verbose && s != nil {
            s.Stop()
        }
        
        if err != nil {
            fmt.Printf("Error uninstalling application: %v\n", err)
        } else {
            fmt.Printf("Application %s uninstalled successfully!\n", appName)
        }
    },
}

func GetCommand() *cobra.Command {
    dependencyCmd.Flags().Bool("verbose", false, "Show Helm output")
    appCmd.Flags().Bool("verbose", false, "Show Helm output")
    uninstallCmd.AddCommand(dependencyCmd)
    uninstallCmd.AddCommand(appCmd)
    return uninstallCmd
}