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
package install

import (
    "fmt"
    "time"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/briandowns/spinner"
    "go-cli/internal/deploy"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
    Use:   "install",
    Short: "Install resources to the cluster",
    Long:  `Install various resources like dependencies to the cluster.`,
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
    },
}

// dependencyCmd represents the dependency subcommand
var dependencyCmd = &cobra.Command{
    Use:   "dependency [dependency-name]",
    Short: "Install a specific dependency",
    Long:  `Install a specific dependency like PostgreSQL, Redis, etc. to the cluster.`,
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
            s.Suffix = fmt.Sprintf(" Installing dependency %s...", depName)
            s.Start()
        }
        
        err := deploy.InstallDependency(depName, depConfig, verbose)
        
        if !verbose && s != nil {
            s.Stop()
        }
        
        if err != nil {
            fmt.Printf("Error installing dependency '%s': %v\n", depName, err)
        } else {
            fmt.Printf("Dependency '%s' installed successfully!\n", depName)
        }
    },
}

// appCmd represents the app subcommand
var appCmd = &cobra.Command{
    Use:   "app [app-name]",
    Short: "Install an application",
    Long:  `Install a specific application to the cluster.`,
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
            s.Suffix = fmt.Sprintf(" Installing application %s...", appName)
            s.Start()
        }
        
        err := deploy.InstallApp(config, appName, verbose)
        
        if !verbose && s != nil {
            s.Stop()
        }
        
        if err != nil {
            fmt.Printf("Error installing application: %v\n", err)
        } else {
            fmt.Printf("Application %s installed successfully!\n", appName)
        }
    },
}

func GetCommand() *cobra.Command {
    dependencyCmd.Flags().Bool("verbose", false, "Show Helm output")
    appCmd.Flags().Bool("verbose", false, "Show Helm output")
    installCmd.AddCommand(dependencyCmd)
    installCmd.AddCommand(appCmd)
    return installCmd
}