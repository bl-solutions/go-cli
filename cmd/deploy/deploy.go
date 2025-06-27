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
    "time"

    "github.com/spf13/cobra"
    "github.com/briandowns/spinner"
    "go-cli/internal/deploy"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Deploy resources to the cluster",
    Long:  `Deploy various resources and dependencies to the cluster.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("deploy called")
    },
}

// dependenciesCmd represents the dependencies subcommand
var dependenciesCmd = &cobra.Command{
    Use:   "dependencies",
    Short: "Deploy application dependencies",
    Long:  `Deploy application dependencies like PostgreSQL, Redis, etc. to the cluster.`,
    Run: func(cmd *cobra.Command, args []string) {
        optional, _ := cmd.Flags().GetBool("optional")
        
        var message string
        if optional {
            message = " Deploying dependencies (including optional)..."
        } else {
            message = " Deploying dependencies..."
        }
        
        s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
        s.Suffix = message
        s.Start()
        
        err := deploy.DeployDependencies(optional)
        
        s.Stop()
        if err != nil {
            fmt.Printf("Error deploying dependencies: %v\n", err)
        } else {
            if optional {
                fmt.Println("Dependencies (including optional) deployed successfully!")
            } else {
                fmt.Println("Dependencies deployed successfully!")
            }
        }
    },
}

// appCmd represents the app subcommand
var appCmd = &cobra.Command{
    Use:   "app [app-name]",
    Short: "Deploy an application",
    Long:  `Deploy a specific application to the cluster.`,
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        appName := args[0]
        
        s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
        s.Suffix = fmt.Sprintf(" Deploying application %s...", appName)
        s.Start()
        
        err := deploy.DeployApp(appName)
        
        s.Stop()
        if err != nil {
            fmt.Printf("Error deploying application: %v\n", err)
        } else {
            fmt.Printf("Application %s deployed successfully!\n", appName)
        }
    },
}

func GetCommand() *cobra.Command {
    dependenciesCmd.Flags().Bool("optional", false, "Deploy optional dependencies")
    deployCmd.AddCommand(dependenciesCmd)
    deployCmd.AddCommand(appCmd)
    return deployCmd
}