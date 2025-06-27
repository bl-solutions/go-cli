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

    "github.com/spf13/cobra"
    "go-cli/internal/cluster"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
    Use:   "cluster",
    Short: "A brief description of your command",
    Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("cluster called")
    },
}

// createCmd represents the create subcommand
var createCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new cluster",
    Long:  `Create a new cluster with the specified configuration.`,
    Run: func(cmd *cobra.Command, args []string) {
        err := cluster.Create()
        if err != nil {
            fmt.Printf("Error creating cluster: %v\n", err)
        }
    },
}

// deleteCmd represents the delete subcommand
var deleteCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete a cluster",
    Long:  `Delete an existing cluster.`,
    Run: func(cmd *cobra.Command, args []string) {
        err := cluster.Delete()
        if err != nil {
            fmt.Printf("Error deleting cluster: %v\n", err)
        }
    },
}

// startCmd represents the start subcommand
var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start a cluster",
    Long:  `Start an existing cluster.`,
    Run: func(cmd *cobra.Command, args []string) {
        err := cluster.Start()
        if err != nil {
            fmt.Printf("Error starting cluster: %v\n", err)
        }
    },
}

// stopCmd represents the stop subcommand
var stopCmd = &cobra.Command{
    Use:   "stop",
    Short: "Stop a cluster",
    Long:  `Stop a running cluster.`,
    Run: func(cmd *cobra.Command, args []string) {
        err := cluster.Stop()
        if err != nil {
            fmt.Printf("Error stopping cluster: %v\n", err)
        }
    },
}

func GetCommand() *cobra.Command {
    clusterCmd.AddCommand(createCmd)
    clusterCmd.AddCommand(deleteCmd)
    clusterCmd.AddCommand(startCmd)
    clusterCmd.AddCommand(stopCmd)
    return clusterCmd
}
